const form = document.getElementById('studentForm');
const resultContainer = document.getElementById('result');

form.addEventListener('submit', async e => {
    e.preventDefault();
    const data = new FormData(form);
    const payload = {
        name: data.get('name'),
        wants_ac: data.get('wants_ac') === 'on',
        wants_dining: data.get('wants_dining') === 'on',
        wants_kitchen: data.get('wants_kitchen') === 'on',
        wants_private_bath: data.get('wants_private_bath') === 'on'
    };

    try {
        const res = await fetch('http://3.14.131.58/api/student', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });
        if (!res.ok) throw new Error(`Status ${res.status}`);
        const student = await res.json();
        displayStudent(student);
    } catch (err) {
        resultContainer.textContent = 'Error: ' + err.message;
    }
});

function displayStudent(s) {
    resultContainer.innerHTML = `
        <div class="student-card">
            <h2>${s.name} (ID: ${s.student_id})</h2>
            <ul>
                <li>Wants AC: ${s.wants_ac ? '✅' : '❌'}</li>
                <li>Wants Dining: ${s.wants_dining ? '✅' : '❌'}</li>
                <li>Wants Kitchen: ${s.wants_kitchen ? '✅' : '❌'}</li>
                <li>Wants Private Bath: ${s.wants_private_bath ? '✅' : '❌'}</li>
            </ul>
        </div>
    `;
}
