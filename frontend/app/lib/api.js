const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function fetchProblems() {
    const res = await fetch(`${API_URL}/api/problems`);
    if (!res.ok) throw new Error('Failed to fetch problems');
    return res.json();
}

export async function fetchProblem(id) {
    const res = await fetch(`${API_URL}/api/problems/${id}`);
    if (!res.ok) throw new Error('Failed to fetch problem');
    return res.json();
}

export async function createSubmission(userId, problemId, code, language = 'python') {
    const res = await fetch(`${API_URL}/api/submissions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            user_id: userId,
            problem_id: problemId,
            code,
            language,
        }),
    });
    if (!res.ok) throw new Error('Failed to create submission');
    return res.json();
}

export async function fetchSubmission(id) {
    const res = await fetch(`${API_URL}/api/submissions/${id}`);
    if (!res.ok) throw new Error('Failed to fetch submission');
    return res.json();
}

export async function createUser(nickname) {
    const res = await fetch(`${API_URL}/api/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname }),
    });
    if (!res.ok) throw new Error('Failed to create user');
    return res.json();
}

// --- Race API ---

export async function createRace(userId, problemId) {
    const res = await fetch(`${API_URL}/api/races`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, problem_id: problemId }),
    });
    if (!res.ok) throw new Error('Failed to create race');
    return res.json();
}

export async function fetchRace(roomCode) {
    const res = await fetch(`${API_URL}/api/races/${roomCode}`);
    if (!res.ok) throw new Error('Failed to fetch race');
    return res.json();
}

export async function joinRace(roomCode, userId) {
    const res = await fetch(`${API_URL}/api/races/${roomCode}/join`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId }),
    });
    if (!res.ok) throw new Error('Failed to join race');
    return res.json();
}

export async function startRace(roomCode, userId) {
    const res = await fetch(`${API_URL}/api/races/${roomCode}/start`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId }),
    });
    if (!res.ok) throw new Error('Failed to start race');
    return res.json();
}

export async function watchRace(roomCode, userId) {
    const res = await fetch(`${API_URL}/api/races/${roomCode}/watch`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId }),
    });
    if (!res.ok) throw new Error('Failed to join as spectator');
    return res.json();
}
