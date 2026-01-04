'use client';

import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useUser } from './hooks/useUser';
import AuthScreen from './components/AuthScreen';
import DashboardLayout from './components/DashboardLayout';
import ContestCard from './components/ContestCard';
import { fetchProblems } from './lib/api';
import styles from './page.module.css';
import Link from 'next/link';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function Home() {
  const { t } = useTranslation();
  const { user, loading: userLoading, login, register, logout } = useUser();
  const [problems, setProblems] = useState([]);
  const [races, setRaces] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (user) {
      loadData();
    }
  }, [user]);

  const loadData = async () => {
    try {
      setLoading(true);
      const [problemsData] = await Promise.all([
        fetchProblems(),
        // fetchRaces() - would fetch active races if API exists
      ]);
      setProblems(problemsData || []);
      // Mock races for demo
      setRaces([
        { room_code: 'DEMO1', title: 'Weekly Contest', status: 'waiting', player_count: 12, duration: 30 },
        { room_code: 'DEMO2', title: 'Quick Match', status: 'active', player_count: 4, duration: 15 },
      ]);
    } catch (err) {
      setError(t('home.errorLoading'));
    } finally {
      setLoading(false);
    }
  };

  if (userLoading) {
    return (
      <div className={styles.loadingScreen}>
        <div className={styles.spinner}></div>
      </div>
    );
  }

  if (!user) {
    return <AuthScreen onLogin={login} onRegister={register} />;
  }

  return (
    <DashboardLayout user={user} onLogout={logout}>
      {/* Active Races Section */}
      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>{t('dashboard.activeContests')}</h2>
          <Link href="/race" className={styles.createBtn}>
            {t('dashboard.createRace')}
          </Link>
        </div>
        <div className={styles.contestsRow}>
          {races.length > 0 ? (
            races.map((race) => (
              <ContestCard key={race.room_code} contest={race} />
            ))
          ) : (
            <p className={styles.emptyText}>{t('dashboard.noActiveContests')}</p>
          )}
        </div>
      </section>

      {/* Problems Section */}
      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>{t('dashboard.problems.title')}</h2>
        </div>

        {loading && (
          <div className={styles.loading}>
            <div className={styles.spinner}></div>
            <span>{t('home.loadingProblems')}</span>
          </div>
        )}

        {error && (
          <div className={styles.error}>
            <p>{error}</p>
            <button onClick={loadData} className={styles.retryBtn}>
              {t('common.retry')}
            </button>
          </div>
        )}

        {!loading && !error && problems.length === 0 && (
          <div className={styles.empty}>
            <p>{t('home.noProblems')}</p>
          </div>
        )}

        {!loading && !error && problems.length > 0 && (
          <div className={styles.problemTable}>
            <div className={styles.tableHeader}>
              <span className={styles.colId}>#</span>
              <span className={styles.colTitle}>{t('dashboard.problems.title')}</span>
              <span className={styles.colTime}>{t('dashboard.problems.timeLimit')}</span>
              <span className={styles.colMemory}>{t('dashboard.problems.memoryLimit')}</span>
              <span className={styles.colAction}></span>
            </div>
            {problems.map((problem) => (
              <Link
                href={`/problem/${problem.id}`}
                key={problem.id}
                className={styles.tableRow}
              >
                <span className={styles.colId}>{problem.id}</span>
                <span className={styles.colTitle}>{problem.title}</span>
                <span className={styles.colTime}>{problem.time_limit_ms}ms</span>
                <span className={styles.colMemory}>{problem.memory_limit_mb}MB</span>
                <span className={styles.colAction}>
                  <button className={styles.solveBtn}>
                    {t('dashboard.problems.solveNow')}
                  </button>
                </span>
              </Link>
            ))}
          </div>
        )}
      </section>
    </DashboardLayout>
  );
}
