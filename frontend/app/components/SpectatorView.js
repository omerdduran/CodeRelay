'use client';

import { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import PlayerCodePanel from './PlayerCodePanel';
import styles from './SpectatorView.module.css';

export default function SpectatorView({ roomCode, players, currentUserId }) {
    const { t } = useTranslation();
    const [viewMode, setViewMode] = useState('grid'); // 'grid' or 'focus'
    const [focusedPlayer, setFocusedPlayer] = useState(null);
    const [playerCodes, setPlayerCodes] = useState({});
    const [playerStatuses, setPlayerStatuses] = useState({});
    const wsRef = useRef(null);

    // Filter only actual players (not spectators)
    const activePlayers = players.filter(p => p.role === 'player');

    useEffect(() => {
        // Connect to WebSocket
        const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';
        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
            // Join the room
            ws.send(JSON.stringify({
                type: 'join_room',
                room_code: roomCode,
                user_id: currentUserId
            }));
        };

        ws.onmessage = (event) => {
            try {
                const messages = event.data.split('\n');
                messages.forEach(msgStr => {
                    if (!msgStr) return;
                    const msg = JSON.parse(msgStr);

                    if (msg.type === 'code_update') {
                        // Parse payload if it's a string, otherwise use directly
                        const payload = typeof msg.payload === 'string' 
                            ? JSON.parse(msg.payload) 
                            : msg.payload;
                        setPlayerCodes(prev => ({
                            ...prev,
                            [payload.user_id]: {
                                code: payload.code,
                                cursor: payload.cursor,
                                nickname: payload.nickname
                            }
                        }));
                    } else if (msg.type === 'player_status') {
                        // Parse payload if it's a string, otherwise use directly
                        const payload = typeof msg.payload === 'string' 
                            ? JSON.parse(msg.payload) 
                            : msg.payload;
                        setPlayerStatuses(prev => ({
                            ...prev,
                            [payload.user_id]: payload.status
                        }));
                    }
                });
            } catch (e) {
                console.error('WS parse error:', e);
            }
        };

        return () => {
            if (ws.readyState === WebSocket.OPEN) {
                ws.close();
            }
        };
    }, [roomCode, currentUserId]);

    const handlePlayerClick = (player) => {
        setFocusedPlayer(player);
        setViewMode('focus');
    };

    const handleBackToGrid = () => {
        setViewMode('grid');
        setFocusedPlayer(null);
    };

    const handleNavigatePlayer = (direction) => {
        const currentIndex = activePlayers.findIndex(p => p.user_id === focusedPlayer?.user_id);
        let newIndex;
        if (direction === 'prev') {
            newIndex = currentIndex <= 0 ? activePlayers.length - 1 : currentIndex - 1;
        } else {
            newIndex = currentIndex >= activePlayers.length - 1 ? 0 : currentIndex + 1;
        }
        setFocusedPlayer(activePlayers[newIndex]);
    };

    // Calculate grid columns based on player count
    const getGridColumns = () => {
        const count = activePlayers.length;
        if (count <= 2) return 2;
        if (count <= 4) return 2;
        if (count <= 6) return 3;
        return 4;
    };

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <div className={styles.headerLeft}>
                    <span className={styles.spectatorBadge}>👁️ {t('race.spectator')}</span>
                    <span className={styles.roomCode}>{roomCode}</span>
                </div>
                <div className={styles.headerRight}>
                    <span className={styles.playerCount}>
                        {activePlayers.length} {t('dashboard.contest.players')}
                    </span>
                    <div className={styles.viewToggle}>
                        <button
                            className={`${styles.toggleBtn} ${viewMode === 'grid' ? styles.active : ''}`}
                            onClick={() => setViewMode('grid')}
                        >
                            ⊞ {t('race.spectatorView.grid')}
                        </button>
                        <button
                            className={`${styles.toggleBtn} ${viewMode === 'focus' ? styles.active : ''}`}
                            onClick={() => focusedPlayer ? setViewMode('focus') : null}
                            disabled={!focusedPlayer}
                        >
                            ⊡ {t('race.spectatorView.focus')}
                        </button>
                    </div>
                </div>
            </div>

            {viewMode === 'grid' ? (
                <div
                    className={styles.grid}
                    style={{ gridTemplateColumns: `repeat(${getGridColumns()}, 1fr)` }}
                >
                    {activePlayers.map(player => (
                        <PlayerCodePanel
                            key={player.user_id}
                            player={player}
                            code={playerCodes[player.user_id]?.code || `# ${t('race.spectatorView.waitingForCode')}`}
                            status={playerStatuses[player.user_id] || 'idle'}
                            isCompact={true}
                            onClick={() => handlePlayerClick(player)}
                        />
                    ))}
                </div>
            ) : (
                <div className={styles.focusView}>
                    <div className={styles.focusHeader}>
                        <button className={styles.navBtn} onClick={() => handleNavigatePlayer('prev')}>
                            ◀ {t('race.spectatorView.prev')}
                        </button>
                        <span className={styles.focusTitle}>
                            👁️ {t('race.spectatorView.watching')}: {focusedPlayer?.nickname}
                        </span>
                        <button className={styles.navBtn} onClick={() => handleNavigatePlayer('next')}>
                            {t('race.spectatorView.next')} ▶
                        </button>
                    </div>
                    {focusedPlayer && (
                        <PlayerCodePanel
                            player={focusedPlayer}
                            code={playerCodes[focusedPlayer.user_id]?.code || `# ${t('race.spectatorView.waitingForCode')}`}
                            status={playerStatuses[focusedPlayer.user_id] || 'idle'}
                            isCompact={false}
                        />
                    )}
                    <button className={styles.backToGridBtn} onClick={handleBackToGrid}>
                        ⊞ {t('race.spectatorView.backToGrid')}
                    </button>
                </div>
            )}
        </div>
    );
}
