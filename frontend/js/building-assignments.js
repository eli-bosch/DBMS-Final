const form = document.getElementById('lookupForm');
const resultContainer = document.getElementById('result');

form.addEventListener('submit', async e => {
    e.preventDefault();
    resultContainer.innerHTML = '';

    const building = form.building_name.value.trim();
    if (!building) return;

    try {
        const res = await fetch(`/api/assignment/${encodeURIComponent(building)}`);
        if (!res.ok) {
            const errText = await res.text();
            throw new Error(errText || `Status ${res.status}`);
        }
        const list = await res.json();
        if (list.length === 0) {
            resultContainer.innerHTML = `<div class="error">No assignments found for “${building}”.</div>`;
        } else {
            displayAssignments(list);
        }
    } catch (err) {
        resultContainer.innerHTML = `<div class="error">${err.message}</div>`;
    }
});

function displayAssignments(assignments) {
    resultContainer.innerHTML = assignments.map(a => `
        <div class="assignment-item">
            <h2>${a.student_name} <span class="sub">(#${a.student_id})</span></h2>
            <p><strong>Room:</strong> Building ${a.building_id}, Room ${a.room_number}</p>
        </div>
    `).join('');
}
