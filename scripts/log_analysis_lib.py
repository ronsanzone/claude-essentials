#!/usr/bin/env python3
"""
Reusable log analysis functions for Claude Code.
Can be imported by both sub-agents and skills.
"""

import json
import sys
from pathlib import Path
from collections import defaultdict, Counter
from datetime import datetime
from typing import Dict, List, Optional, Any
import statistics

class LogAnalyzer:
    """Base class for log analysis operations."""

    def __init__(self, file_path: str):
        self.file_path = Path(file_path)
        self.validate_file()

    def validate_file(self):
        """Ensure file exists and is readable."""
        if not self.file_path.exists():
            raise FileNotFoundError(f"File not found: {self.file_path}")

        # Check file size for memory warnings
        size_mb = self.file_path.stat().st_size / (1024 * 1024)
        if size_mb > 100:
            print(f"Warning: Large file ({size_mb:.1f}MB). Using streaming mode.",
                  file=sys.stderr)

    def stream_logs(self):
        """Generator to stream logs line by line, handling Splunk export format."""
        with open(self.file_path, 'r') as f:
            for line_num, line in enumerate(f, 1):
                if line_num % 10000 == 0:
                    print(f"Progress: {line_num} lines processed...",
                          file=sys.stderr)

                try:
                    data = json.loads(line.strip())

                    # Handle Splunk export format
                    if 'result' in data and '_raw' in data['result']:
                        # Extract the actual log from Splunk wrapper
                        raw_log = json.loads(data['result']['_raw'])

                        # Normalize field names for consistency
                        # Handle correlationID vs correlationId
                        if 'correlationId' not in raw_log and 'correlationID' in raw_log:
                            raw_log['correlationId'] = raw_log['correlationID']

                        # Extract service name from various sources
                        if 'service' not in raw_log:
                            if 'logger' in raw_log:
                                raw_log['service'] = raw_log['logger']
                            elif 'kube' in raw_log and 'labels' in raw_log['kube']:
                                raw_log['service'] = raw_log['kube']['labels'].get('xgen_app', 'unknown')
                            else:
                                raw_log['service'] = 'unknown'

                        # Normalize timestamp field
                        if 'timestamp' not in raw_log and 'time' in raw_log:
                            raw_log['timestamp'] = raw_log['time']

                        # Normalize message field
                        if 'message' not in raw_log and 'msg' in raw_log:
                            raw_log['message'] = raw_log['msg']

                        # Detect errors even if level is "info"
                        if raw_log.get('level') == 'info':
                            # Check for error indicators
                            if any(field in raw_log for field in ['error', 'grpc.error_message', 'status.message']):
                                raw_log['level'] = 'ERROR'
                            elif 'grpc.success' in raw_log and raw_log['grpc.success'] == 'false':
                                raw_log['level'] = 'ERROR'

                        yield raw_log
                    else:
                        # Fallback to direct parsing if not Splunk format
                        yield data
                except (json.JSONDecodeError, KeyError, TypeError):
                    continue

    def parse_timestamp(self, ts_str: str) -> Optional[datetime]:
        """Parse various timestamp formats."""
        if not ts_str:
            return None

        formats = [
            '%Y-%m-%dT%H:%M:%S.%fZ',
            '%Y-%m-%d %H:%M:%S',
            '%Y-%m-%dT%H:%M:%S%z',
            '%Y-%m-%dT%H:%M:%S.%f%z'
        ]

        # Handle truncation for parsing
        ts_clean = ts_str[:26] if len(ts_str) > 26 else ts_str

        for fmt in formats:
            try:
                return datetime.strptime(ts_clean, fmt)
            except (ValueError, TypeError):
                continue
        return None

class CorrelationAnalyzer(LogAnalyzer):
    """Analyze logs by correlation ID."""

    def analyze(self, correlation_id: Optional[str] = None) -> Dict[str, Any]:
        """Perform correlation analysis with streaming."""
        correlations = defaultdict(list)
        error_patterns = Counter()
        service_errors = defaultdict(Counter)
        first_seen = {}
        last_seen = {}
        line_count = 0
        error_count = 0

        for log in self.stream_logs():
            line_count += 1
            corr_id = log.get('correlationId', 'unknown')

            # Filter if specific correlation requested
            if correlation_id and corr_id != correlation_id:
                continue

            # Track timeline
            timestamp = self.parse_timestamp(log.get('timestamp', ''))
            if timestamp and corr_id != 'unknown':
                if corr_id not in first_seen:
                    first_seen[corr_id] = timestamp
                last_seen[corr_id] = timestamp

            if log.get('level') in ['ERROR', 'FATAL']:
                error_count += 1

                # Track limited data per correlation (memory management)
                if len(correlations[corr_id]) < 100:
                    correlations[corr_id].append({
                        'timestamp': log.get('timestamp'),
                        'service': log.get('service', 'unknown'),
                        'caller': log.get('caller', ''),
                        'message': log.get('message', '')[:200],
                        'level': log.get('level')
                    })

                # Track error patterns
                service = log.get('service', 'unknown')
                caller = log.get('caller', '')
                message = log.get('message', '')[:50]

                error_key = f"{service}:{caller}:{message}"
                error_patterns[error_key] += 1
                service_errors[service][message[:100]] += 1

        # Calculate statistics
        correlation_stats = {}
        for corr_id, events in correlations.items():
            if events and corr_id in first_seen and corr_id in last_seen:
                duration = (last_seen[corr_id] - first_seen[corr_id]).total_seconds()
                correlation_stats[corr_id] = {
                    'error_count': len(events),
                    'duration_seconds': duration,
                    'services_affected': list(set(e['service'] for e in events)),
                    'first_error': events[0],
                    'last_error': events[-1]
                }

        # Find top patterns
        top_errors = error_patterns.most_common(10)
        problem_correlations = sorted(
            correlation_stats.items(),
            key=lambda x: x[1]['error_count'],
            reverse=True
        )[:10]

        # Time range
        time_range = None
        if first_seen and last_seen:
            all_first = min(first_seen.values())
            all_last = max(last_seen.values())
            time_range = {
                'start': all_first.isoformat(),
                'end': all_last.isoformat()
            }

        return {
            'summary': {
                'total_lines': line_count,
                'total_errors': error_count,
                'unique_correlations': len(correlations),
                'time_range': time_range
            },
            'top_error_patterns': [
                {'pattern': k, 'count': v} for k, v in top_errors
            ],
            'problem_correlations': [
                {'correlation_id': k, **v} for k, v in problem_correlations
            ],
            'service_errors': {
                service: dict(errors.most_common(5))
                for service, errors in service_errors.items()
            }
        }

class CascadeDetector(LogAnalyzer):
    """Detect cascading failures in logs."""

    def detect(self, window_seconds: int = 5) -> Dict[str, Any]:
        """Detect cascade patterns within time windows."""
        from datetime import timedelta

        time_buckets = defaultdict(list)

        for log in self.stream_logs():
            if log.get('level') not in ['ERROR', 'FATAL']:
                continue

            timestamp = self.parse_timestamp(log.get('timestamp', ''))
            if not timestamp:
                continue

            # Group by time window
            bucket_key = timestamp.replace(microsecond=0)
            time_buckets[bucket_key].append({
                'timestamp': timestamp,
                'service': log.get('service', 'unknown'),
                'correlation_id': log.get('correlationId', ''),
                'caller': log.get('caller', ''),
                'message': log.get('message', '')[:100]
            })

        # Find cascades
        cascades = []
        for bucket_time, events in sorted(time_buckets.items()):
            if len(events) >= 3:
                # Sort events by exact timestamp
                events.sort(key=lambda x: x['timestamp'])

                # Get unique services while preserving order
                services = []
                for event in events:
                    if event['service'] not in services:
                        services.append(event['service'])

                # Cascade = 3+ different services failing
                if len(services) >= 3:
                    duration = (events[-1]['timestamp'] - events[0]['timestamp']).total_seconds()

                    cascades.append({
                        'start_time': events[0]['timestamp'].isoformat(),
                        'duration_seconds': duration,
                        'services_affected': services[:10],
                        'event_count': len(events),
                        'correlations_affected': list(set(e['correlation_id'] for e in events if e['correlation_id']))[:20],
                        'pattern': ' -> '.join(services[:5])
                    })

        # Sort by event count (severity)
        cascades.sort(key=lambda x: x['event_count'], reverse=True)

        return {
            'cascades_detected': len(cascades),
            'cascade_patterns': cascades[:10],
            'analysis_window_seconds': window_seconds
        }

class PerformanceProfiler(LogAnalyzer):
    """Profile request performance and identify bottlenecks."""

    def profile(self) -> Dict[str, Any]:
        """Analyze performance metrics from logs."""

        operation_times = defaultdict(list)
        correlation_times = {}
        slow_operations = []

        for log in self.stream_logs():
            correlation_id = log.get('correlationId')
            if not correlation_id:
                continue

            # Track correlation lifecycle
            timestamp = self.parse_timestamp(log.get('timestamp', ''))
            if timestamp:
                if correlation_id not in correlation_times:
                    correlation_times[correlation_id] = {
                        'start': timestamp,
                        'end': timestamp
                    }
                else:
                    correlation_times[correlation_id]['end'] = timestamp

            # Extract performance metrics
            if 'duration_ms' in log:
                operation = f"{log.get('service', 'unknown')}:{log.get('operation', 'unknown')}"
                duration = log['duration_ms']
                operation_times[operation].append(duration)

                if duration > 1000:  # Slow operation (>1 second)
                    slow_operations.append({
                        'correlation_id': correlation_id,
                        'operation': operation,
                        'duration_ms': duration,
                        'timestamp': timestamp.isoformat() if timestamp else None,
                        'caller': log.get('caller', '')
                    })

        # Calculate statistics
        operation_stats = {}
        for operation, times in operation_times.items():
            if times:
                stats = {
                    'count': len(times),
                    'mean_ms': statistics.mean(times),
                    'median_ms': statistics.median(times),
                    'max_ms': max(times),
                    'min_ms': min(times)
                }

                # Calculate p95 if enough data
                if len(times) > 20:
                    stats['p95_ms'] = statistics.quantiles(times, n=20)[18]
                else:
                    stats['p95_ms'] = max(times)

                operation_stats[operation] = stats

        # Calculate request durations
        request_durations = []
        for corr_id, times in correlation_times.items():
            if times['start'] and times['end'] and times['start'] != times['end']:
                duration_ms = (times['end'] - times['start']).total_seconds() * 1000
                request_durations.append({
                    'correlation_id': corr_id,
                    'duration_ms': duration_ms
                })

        # Sort results
        slow_operations.sort(key=lambda x: x['duration_ms'], reverse=True)
        request_durations.sort(key=lambda x: x['duration_ms'], reverse=True)

        # Sort operations by p95
        sorted_ops = dict(sorted(
            operation_stats.items(),
            key=lambda x: x[1].get('p95_ms', 0),
            reverse=True
        )[:20])

        return {
            'operation_statistics': sorted_ops,
            'slowest_operations': slow_operations[:20],
            'slowest_requests': request_durations[:20],
            'summary': {
                'total_operations_profiled': sum(s['count'] for s in operation_stats.values()),
                'unique_operations': len(operation_stats),
                'requests_analyzed': len(correlation_times)
            }
        }

class ErrorAnalyzer(LogAnalyzer):
    """Quick error analysis by service or pattern."""

    def analyze_errors(self, service: Optional[str] = None) -> Dict[str, Any]:
        """Analyze error patterns, optionally filtered by service."""

        errors_by_level = Counter()
        errors_by_service = defaultdict(Counter)
        error_messages = Counter()
        error_timeline = defaultdict(int)
        caller_errors = Counter()
        error_types = Counter()
        grpc_errors = Counter()

        for log in self.stream_logs():
            # Extract error message from various fields
            error_msg = None
            is_error = False

            # Check traditional level field
            level = log.get('level', '')
            if level in ['ERROR', 'FATAL', 'WARNING']:
                is_error = True

            # Check for error fields even in info logs
            if 'error' in log:
                is_error = True
                error_msg = log['error']
            elif 'grpc.error_message' in log:
                is_error = True
                error_msg = log['grpc.error_message']
            elif 'status.message' in log:
                is_error = True
                error_msg = log['status.message']
            elif log.get('grpc.success') == 'false':
                is_error = True
                error_msg = log.get('message', log.get('msg', 'grpc failure'))

            if not is_error:
                continue

            # Determine actual level
            if level not in ['ERROR', 'FATAL', 'WARNING']:
                level = 'ERROR'  # Default error level for detected errors

            log_service = log.get('service', 'unknown')

            # Filter by service if specified
            if service and log_service != service:
                continue

            errors_by_level[level] += 1
            errors_by_service[log_service][level] += 1

            # Track error messages
            if error_msg:
                # Clean up and truncate error message
                clean_msg = error_msg.split(': ')[-1][:100]  # Get last part after colon
                error_messages[clean_msg] += 1

                # Categorize error types
                if 'context canceled' in error_msg.lower():
                    error_types['context_canceled'] += 1
                elif 'timeout' in error_msg.lower():
                    error_types['timeout'] += 1
                elif 's3' in error_msg.lower():
                    error_types['s3_error'] += 1
                elif 'unavailable' in error_msg.lower():
                    error_types['service_unavailable'] += 1
                elif 'permission' in error_msg.lower() or 'unauthorized' in error_msg.lower():
                    error_types['auth_error'] += 1
                else:
                    error_types['other'] += 1

            # Track gRPC errors specifically
            if 'grpc.code' in log and log['grpc.code'] != 'OK':
                grpc_errors[log['grpc.code']] += 1

            # Use message field as fallback
            if not error_msg:
                message = log.get('message', '')[:100]
                if message:
                    error_messages[message] += 1

            # Track callers
            caller = log.get('caller', '')
            if caller:
                caller_errors[caller] += 1

            # Timeline (5-minute buckets)
            timestamp = self.parse_timestamp(log.get('timestamp', ''))
            if timestamp:
                bucket = timestamp.strftime('%Y-%m-%d %H:%M')[:-1] + '0'
                error_timeline[bucket] += 1

        # Find error spikes
        error_spikes = sorted(
            [(t, count) for t, count in error_timeline.items() if count > 10],
            key=lambda x: x[1],
            reverse=True
        )[:10]

        return {
            'summary': {
                'total_errors': sum(errors_by_level.values()),
                'errors_by_level': dict(errors_by_level),
                'services_with_errors': len(errors_by_service),
                'error_types': dict(error_types),
                'grpc_errors': dict(grpc_errors)
            },
            'top_error_messages': dict(error_messages.most_common(15)),
            'errors_by_service': {
                svc: dict(counts)
                for svc, counts in list(errors_by_service.items())[:10]
            },
            'top_error_locations': dict(caller_errors.most_common(15)),
            'error_spikes': [
                {'time': t, 'count': c} for t, c in error_spikes
            ]
        }

# CLI Interface
def main():
    """Command-line interface for direct execution."""
    import argparse

    parser = argparse.ArgumentParser(description='Analyze Splunk JSON logs')
    parser.add_argument('command',
                       choices=['correlate', 'cascade', 'profile', 'errors'],
                       help='Analysis command to run')
    parser.add_argument('file', help='Log file path')
    parser.add_argument('--correlation-id', help='Specific correlation to analyze')
    parser.add_argument('--window', type=int, default=5,
                       help='Time window for cascade detection (seconds)')
    parser.add_argument('--service', help='Filter errors by service')

    args = parser.parse_args()

    try:
        if args.command == 'correlate':
            analyzer = CorrelationAnalyzer(args.file)
            results = analyzer.analyze(args.correlation_id)
        elif args.command == 'cascade':
            detector = CascadeDetector(args.file)
            results = detector.detect(args.window)
        elif args.command == 'profile':
            profiler = PerformanceProfiler(args.file)
            results = profiler.profile()
        elif args.command == 'errors':
            analyzer = ErrorAnalyzer(args.file)
            results = analyzer.analyze_errors(args.service)
        else:
            results = {'error': f'Command {args.command} not implemented'}

        print(json.dumps(results, indent=2, default=str))

    except Exception as e:
        print(json.dumps({'error': str(e)}, indent=2), file=sys.stderr)
        sys.exit(1)

if __name__ == '__main__':
    main()
