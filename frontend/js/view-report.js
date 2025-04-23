const btn = document.getElementById('loadReport');
const tbody = document.querySelector('#reportTable tbody');

btn.addEventListener('click', async () => {
    tbody.innerHTML = '';
    try {
        const res = await fetch('http://3.14.131.58/api/building/report');
        if (!res.ok) throw new Error(`Status ${res.status}`);
        let data = await res.json();
        data = data.filter(r => r.building_id === 1);
        data.forEach(row => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${row.building_name}</td>
                <td>${row.total_rooms ?? '—'}</td>
                <td>${row.total_bedrooms ?? '—'}</td>
                <td>${row.rooms_with_availability ?? '—'}</td>
                <td>${row.available_bedrooms}</td>
            `;
            tbody.appendChild(tr);
        });
    } catch (err) {
        const tr = document.createElement('tr');
        tr.innerHTML = `<td colspan="5" style="color:#C00;">Error: ${err.message}</td>`;
        tbody.appendChild(tr);
    }
});
