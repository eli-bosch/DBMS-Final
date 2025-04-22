const form = document.getElementById('prefForm');
const result = document.getElementById('result');

form.addEventListener('submit', async e => {
    e.preventDefault();
    result.innerHTML = '';

    const id = form.student_id.value.trim();
    if (!id) return;

    try {
        const res = await fetch(`http://3.14.131.58/api/preference/${id}`);
        if (!res.ok) {
            const txt = await res.text();
            throw new Error(txt || `Status ${res.status}`);
        }
        const rooms = await res.json();
        if (rooms.length === 0) {
            result.innerHTML = `<div class="error">No rooms match preferences.</div>`;
            return;
        }
        rooms.forEach(r => {
            const div = document.createElement('div');
            div.className = 'room-item';
            div.innerHTML = `
                <p><strong>Building:</strong> ${r.building_id}</p>
                <p><strong>Room:</strong> ${r.room_number}</p>
                <p><strong>Bedrooms:</strong> ${r.num_bedroom}</p>
                <p><strong>Private Baths:</strong> ${r.private_bathrooms}</p>
                <p><strong>Kitchen:</strong> ${r.has_kitchen ? '✅' : '❌'}</p>
            `;
            result.appendChild(div);
        });
    } catch (err) {
        result.innerHTML = `<div class="error">Error: ${err.message}</div>`;
    }
});
