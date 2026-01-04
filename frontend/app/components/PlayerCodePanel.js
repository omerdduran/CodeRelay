'use client';

import { useMemo } from 'react';
import styles from './PlayerCodePanel.module.css';

const STATUS_COLORS = {
    typing: '#22c55e',
    idle: '#6b7280',
    submitting: '#f59e0b',
    finished: '#8b5cf6',
};

const STATUS_LABELS = {
    typing: '⌨️ Typing',
    idle: '💤 Idle',
    submitting: '🚀 Submitting',
    finished: '✅ Done',
};

export default function PlayerCodePanel({ player, code, status, isCompact, onClick }) {
    // Truncate code for compact view
    const displayCode = useMemo(() => {
        if (!isCompact) return code;
        const lines = code.split('\n');
        if (lines.length > 15) {
            return lines.slice(0, 15).join('\n') + '\n...';
        }
        return code;
    }, [code, isCompact]);

    // Count lines
    const lineCount = code.split('\n').length;

    return (
        <div
            className={`${styles.panel} ${isCompact ? styles.compact : styles.full}`}
            onClick={onClick}
            style={{ cursor: isCompact ? 'pointer' : 'default' }}
        >
            <div className={styles.header}>
                <div className={styles.playerInfo}>
                    <div className={styles.avatar}>
                        {player.nickname?.charAt(0).toUpperCase() || '?'}
                    </div>
                    <span className={styles.nickname}>{player.nickname}</span>
                </div>
                <div className={styles.status} style={{ color: STATUS_COLORS[status] }}>
                    {STATUS_LABELS[status] || status}
                </div>
            </div>
            <div className={styles.codeContainer}>
                <div className={styles.lineNumbers}>
                    {displayCode.split('\n').map((_, i) => (
                        <div key={i} className={styles.lineNumber}>{i + 1}</div>
                    ))}
                </div>
                <pre className={styles.code}>{displayCode}</pre>
            </div>
            {isCompact && (
                <div className={styles.footer}>
                    <span className={styles.lineCount}>{lineCount} lines</span>
                    <span className={styles.clickHint}>Click to focus</span>
                </div>
            )}
        </div>
    );
}
