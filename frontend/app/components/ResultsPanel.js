'use client';

import styles from './ResultsPanel.module.css';

const STATUS_CONFIG = {
    queued: { color: '#94a3b8', icon: '⏳', label: 'Queued' },
    running: { color: '#60a5fa', icon: '⚡', label: 'Running' },
    AC: { color: '#10b981', icon: '✓', label: 'Accepted' },
    WA: { color: '#ef4444', icon: '✗', label: 'Wrong Answer' },
    TLE: { color: '#f59e0b', icon: '⏱', label: 'Time Limit Exceeded' },
};

export default function ResultsPanel({ submission, loading }) {
    if (loading) {
        return (
            <div className={styles.container}>
                <div className={styles.loading}>
                    <div className={styles.spinner}></div>
                    <span>Submitting...</span>
                </div>
            </div>
        );
    }

    if (!submission) {
        return (
            <div className={styles.container}>
                <div className={styles.empty}>
                    <span className={styles.emptyIcon}>📝</span>
                    <p>Submit your code to see results</p>
                </div>
            </div>
        );
    }

    const config = STATUS_CONFIG[submission.status] || STATUS_CONFIG.queued;

    return (
        <div className={styles.container}>
            <div className={styles.result} style={{ borderColor: config.color }}>
                <div className={styles.status} style={{ color: config.color }}>
                    <span className={styles.statusIcon}>{config.icon}</span>
                    <span className={styles.statusLabel}>{config.label}</span>
                </div>

                <div className={styles.details}>
                    <div className={styles.detail}>
                        <span className={styles.detailLabel}>Submission ID</span>
                        <span className={styles.detailValue}>#{submission.id}</span>
                    </div>
                    {submission.runtime_ms && (
                        <div className={styles.detail}>
                            <span className={styles.detailLabel}>Runtime</span>
                            <span className={styles.detailValue}>{submission.runtime_ms}ms</span>
                        </div>
                    )}
                    <div className={styles.detail}>
                        <span className={styles.detailLabel}>Language</span>
                        <span className={styles.detailValue}>{submission.language}</span>
                    </div>
                </div>
            </div>
        </div>
    );
}
