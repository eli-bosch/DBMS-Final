const btn = document.getElementById('loadBtn');
const result = document.getElementById('result');

btn.addEventListener('click', async () => {
    result.innerHTML = '';
    try {
        const res = await fetch('http://3.14.131.58/api/rooms');
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const rooms = await res.json();
        if (rooms.length === 0) {
            result.textContent = 'No rooms found.';
            return;
        }
        rooms.forEach(r => {
            const div = document.createElement('div');
            div.className = 'room-item';
            div.innerHTML = `
                <p><strong>Building:</strong> ${r.building_id}</p>
                <p><strong>Room:</strong> ${r.room_number}</p>
                <p><strong>Bedrooms:</strong> ${r.num_bedrooms}</p>
            `;
            result.appendChild(div);
        });
    } catch (err) {
        result.textContent = `Error: ${err.message}`;
    }
});
