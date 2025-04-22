const form = document.getElementById('assignmentForm');
const resultContainer = document.getElementById('result');

form.addEventListener('submit', async e => {
    e.preventDefault();
    resultContainer.innerHTML = ''; // clear previous

    const data = new FormData(form);
    const payload = {
        student_id: Number(data.get('student_id')),
        building_id: Number(data.get('building_id')),
        room_number: Number(data.get('room_number'))
    };

    try {
        const res = await fetch('http://3.14.131.58/api/assignment', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            const errText = await res.text();
            throw new Error(errText || `Status ${res.status}`);
        }

        const a = await res.json();
        displayAssignment(a);
    } catch (err) {
        resultContainer.innerHTML = `<div class="error">${err.message}</div>`;
    }
});

function displayAssignment(a) {
    resultContainer.innerHTML = `
        <div class="assignment-card">
            <h2>Assigned Student #${a.student_id}</h2>
            <p><strong>Building ID:</strong> ${a.building_id}</p>
            <p><strong>Room Number:</strong> ${a.room_number}</p>
        </div>
    `;
}
