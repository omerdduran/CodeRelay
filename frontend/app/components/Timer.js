'use client';

import { useEffect, useState } from 'react';
import styles from './Timer.module.css';

export default function Timer({
    durationMinutes = 30,
    startTime = null,
    onTimeUp = () => { }
}) {
    const [timeLeft, setTimeLeft] = useState(durationMinutes * 60);
    const [isRunning, setIsRunning] = useState(false);

    useEffect(() => {
        if (startTime) {
            const start = new Date(startTime).getTime();
            const end = start + durationMinutes * 60 * 1000;
            const now = Date.now();
            const remaining = Math.max(0, Math.floor((end - now) / 1000));
            setTimeLeft(remaining);
            setIsRunning(remaining > 0);
        }
    }, [startTime, durationMinutes]);

    useEffect(() => {
        if (!isRunning) return;

        const interval = setInterval(() => {
            setTimeLeft((prev) => {
                if (prev <= 1) {
                    setIsRunning(false);
                    onTimeUp();
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);

        return () => clearInterval(interval);
    }, [isRunning, onTimeUp]);

    const minutes = Math.floor(timeLeft / 60);
    const seconds = timeLeft % 60;

    const isLow = timeLeft < 60;
    const isCritical = timeLeft < 30;

    return (
        <div className={`${styles.container} ${isLow ? styles.low : ''} ${isCritical ? styles.critical : ''}`}>
            <span className={styles.icon}>⏱</span>
            <span className={styles.time}>
                {String(minutes).padStart(2, '0')}:{String(seconds).padStart(2, '0')}
            </span>
        </div>
    );
}
