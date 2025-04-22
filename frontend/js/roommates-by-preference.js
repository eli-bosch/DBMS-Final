const form = document.getElementById('mateForm');
const result = document.getElementById('result');

form.addEventListener('submit', async e => {
    e.preventDefault();
    result.innerHTML = '';

    const id = form.student_id.value.trim();
    if (!id) return;

    try {
        const res = await fetch(`/api/student/${id}`);
        if (!res.ok) {
            const msg = await res.text();
            throw new Error(msg || `Status ${res.status}`);
        }
        const mates = await res.json();
        if (mates.length === 0) {
            result.innerHTML = `<div class="error">No matching roommates found.</div>`;
            return;
        }
        mates.forEach(m => {
            const div = document.createElement('div');
            div.className = 'mate-item';
            div.innerHTML = `
                <h2>${m.name} (ID: ${m.student_id})</h2>
                <ul>
                    <li>Wants AC: ${m.wants_ac ? '✅' : '❌'}</li>
                    <li>Wants Dining: ${m.wants_dining ? '✅' : '❌'}</li>
                    <li>Wants Kitchen: ${m.wants_kitchen ? '✅' : '❌'}</li>
                    <li>Wants Private Bath: ${m.wants_private_bath ? '✅' : '❌'}</li>
                </ul>
            `;
            result.appendChild(div);
        });
    } catch (err) {
        result.innerHTML = `<div class="error">Error: ${err.message}</div>`;
    }
});
