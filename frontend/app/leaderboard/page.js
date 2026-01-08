'use client';

import { useTranslation } from 'react-i18next';
import { useUser } from '../hooks/useUser';
import AuthScreen from '../components/AuthScreen';
import DashboardLayout from '../components/DashboardLayout';
import Leaderboard from '../components/Leaderboard';

export default function LeaderboardPage() {
  const { t } = useTranslation();
  const { user, loading: userLoading, login, register, logout } = useUser();

  if (userLoading) {
    return (
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          height: '100vh',
          color: 'white',
        }}
      >
        {t('common.loading')}
      </div>
    );
  }

  if (!user) {
    return <AuthScreen onLogin={login} onRegister={register} />;
  }

  return (
    <DashboardLayout user={user} onLogout={logout}>
      <Leaderboard showElo={true} />
    </DashboardLayout>
  );
}

