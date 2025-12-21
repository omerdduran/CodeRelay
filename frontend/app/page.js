'use client';

import { useEffect, useState } from 'react';
import { useUser } from './hooks/useUser';
import NicknameScreen from './components/NicknameScreen';
import { fetchProblems } from './lib/api';
import styles from './page.module.css';
import Link from 'next/link';

export default function Home() {
  const { user, loading: userLoading, login } = useUser();
  const [problems, setProblems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (user) {
      loadProblems();
    }
  }, [user]);

  const loadProblems = async () => {
    try {
      setLoading(true);
      const data = await fetchProblems();
      setProblems(data || []);
    } catch (err) {
      setError('Failed to load problems. Make sure the backend is running.');
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
    return <NicknameScreen onNicknameSet={login} />;
  }

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.logo}>
          <span className={styles.logoIcon}>⚡</span>
          <h1 className={styles.title}>CodeRelay</h1>
        </div>
        <div className={styles.userInfo}>
          <span className={styles.nickname}>{user.nickname}</span>
        </div>
      </header>

      <main className={styles.main}>
        <h2 className={styles.sectionTitle}>Problems</h2>

        {loading && (
          <div className={styles.loading}>
            <div className={styles.spinner}></div>
            <span>Loading problems...</span>
          </div>
        )}

        {error && (
          <div className={styles.error}>
            <p>{error}</p>
            <button onClick={loadProblems} className={styles.retryBtn}>
              Retry
            </button>
          </div>
        )}

        {!loading && !error && problems.length === 0 && (
          <div className={styles.empty}>
            <p>No problems available yet.</p>
          </div>
        )}

        {!loading && !error && problems.length > 0 && (
          <div className={styles.problemList}>
            {problems.map((problem) => (
              <Link
                href={`/problem/${problem.id}`}
                key={problem.id}
                className={styles.problemCard}
              >
                <div className={styles.problemHeader}>
                  <span className={styles.problemId}>#{problem.id}</span>
                  <h3 className={styles.problemTitle}>{problem.title}</h3>
                </div>
                <div className={styles.problemMeta}>
                  <span className={styles.problemLimit}>⏱ {problem.time_limit_ms}ms</span>
                  <span className={styles.problemLimit}>💾 {problem.memory_limit_mb}MB</span>
                </div>
                <div className={styles.problemCta}>
                  Solve →
                </div>
              </Link>
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
