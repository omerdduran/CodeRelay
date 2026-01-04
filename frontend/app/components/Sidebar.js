'use client';

import { useTranslation } from 'react-i18next';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import styles from './Sidebar.module.css';

const MENU_ITEMS = [
    { key: 'practice', icon: '💻', href: '/' },
    { key: 'races', icon: '🏁', href: '/race' },
    { key: 'leaderboard', icon: '🏆', href: '/leaderboard' },
];

export default function Sidebar({ user, onLogout }) {
    const { t } = useTranslation();
    const pathname = usePathname();

    return (
        <aside className={styles.sidebar}>
            <div className={styles.logo}>
                <span className={styles.logoIcon}>⚡</span>
                <span className={styles.logoText}>CodeRelay</span>
            </div>

            <nav className={styles.nav}>
                {MENU_ITEMS.map((item) => (
                    <Link
                        key={item.key}
                        href={item.href}
                        className={`${styles.navItem} ${pathname === item.href ? styles.active : ''}`}
                    >
                        <span className={styles.navIcon}>{item.icon}</span>
                        <span className={styles.navLabel}>{t(`dashboard.menu.${item.key}`)}</span>
                    </Link>
                ))}
            </nav>

            <div className={styles.userSection}>
                <div className={styles.userCard}>
                    <div className={styles.userAvatar}>
                        {user?.nickname?.charAt(0).toUpperCase() || '?'}
                    </div>
                    <div className={styles.userInfo}>
                        <span className={styles.userName}>{user?.nickname}</span>
                        <button onClick={onLogout} className={styles.logoutBtn}>
                            {t('common.logout')}
                        </button>
                    </div>
                </div>
            </div>
        </aside>
    );
}
