'use client';

import { useTranslation } from 'react-i18next';
import styles from './ResultsPanel.module.css';

export default function ResultsPanel({ submission, loading }) {
    const { t } = useTranslation();

    const STATUS_CONFIG = {
        queued: { color: '#94a3b8', icon: '⏳', label: t('results.status.queued') },
        running: { color: '#60a5fa', icon: '⚡', label: t('results.status.running') },
        AC: { color: '#10b981', icon: '✓', label: t('results.status.AC') },
        WA: { color: '#ef4444', icon: '✗', label: t('results.status.WA') },
        TLE: { color: '#f59e0b', icon: '⏱', label: t('results.status.TLE') },
    };

    if (loading) {
        return (
            <div className={styles.container}>
                <div className={styles.loading}>
                    <div className={styles.spinner}></div>
                    <span>{t('results.submitting')}</span>
                </div>
            </div>
        );
    }

    if (!submission) {
        return (
            <div className={styles.container}>
                <div className={styles.empty}>
                    <span className={styles.emptyIcon}>📝</span>
                    <p>{t('results.emptyMessage')}</p>
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
                        <span className={styles.detailLabel}>{t('results.submissionId')}</span>
                        <span className={styles.detailValue}>#{submission.id}</span>
                    </div>
                    {submission.runtime_ms && (
                        <div className={styles.detail}>
                            <span className={styles.detailLabel}>{t('results.runtime')}</span>
                            <span className={styles.detailValue}>{submission.runtime_ms}ms</span>
                        </div>
                    )}
                    <div className={styles.detail}>
                        <span className={styles.detailLabel}>{t('results.language')}</span>
                        <span className={styles.detailValue}>{submission.language}</span>
                    </div>
                </div>
            </div>
        </div>
    );
}
