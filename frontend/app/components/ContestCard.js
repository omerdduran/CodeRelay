'use client';

import { useTranslation } from 'react-i18next';
import Link from 'next/link';
import styles from './ContestCard.module.css';

export default function ContestCard({ contest }) {
    const { t } = useTranslation();

    const statusColors = {
        waiting: '#f59e0b',
        active: '#10b981',
        finished: '#6b7280',
    };

    const statusLabels = {
        waiting: t('dashboard.contest.waiting'),
        active: t('dashboard.contest.active'),
        finished: t('dashboard.contest.finished'),
    };

    return (
        <Link href={`/race/${contest.room_code}`} className={styles.card}>
            <div className={styles.header}>
                <span
                    className={styles.status}
                    style={{ backgroundColor: statusColors[contest.status] || statusColors.waiting }}
                >
                    {statusLabels[contest.status] || contest.status}
                </span>
            </div>
            <h3 className={styles.title}>{contest.title || `Race #${contest.room_code}`}</h3>
            <div className={styles.meta}>
                <div className={styles.metaItem}>
                    <span className={styles.metaIcon}>👥</span>
                    <span>{contest.player_count || 0} {t('dashboard.contest.players')}</span>
                </div>
                <div className={styles.metaItem}>
                    <span className={styles.metaIcon}>⏱️</span>
                    <span>{contest.duration || '30'} min</span>
                </div>
            </div>
            <div className={styles.cta}>
                {contest.status === 'waiting' ? t('dashboard.contest.join') : t('dashboard.contest.view')}
            </div>
        </Link>
    );
}
