'use client';

import Sidebar from './Sidebar';
import LanguageSwitcher from './LanguageSwitcher';
import styles from './DashboardLayout.module.css';

export default function DashboardLayout({ user, onLogout, children }) {
    return (
        <div className={styles.layout}>
            <Sidebar user={user} onLogout={onLogout} />
            <div className={styles.content}>
                <header className={styles.header}>
                    <div className={styles.headerLeft}>
                        {/* Breadcrumb or page title can go here */}
                    </div>
                    <div className={styles.headerRight}>
                        <LanguageSwitcher />
                    </div>
                </header>
                <main className={styles.main}>
                    {children}
                </main>
            </div>
        </div>
    );
}
