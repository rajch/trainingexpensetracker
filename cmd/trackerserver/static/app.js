const apiRoot = '/api';

document.addEventListener('DOMContentLoaded', () => {
    const expenseList = document.getElementById('expense-list');
    const modal = document.getElementById('modal');
    const addBtn = document.getElementById('add-btn');
    const closeBtn = document.querySelector('.close');
    const expenseForm = document.getElementById('expense-form');
    const modalTitle = document.getElementById('modal-title');
    const refreshBanner = document.getElementById('refresh-banner');

    // Load expenses on start
    loadExpenses();

    // // WebSocket logic
    // const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    // const wsUrl = `${protocol}//${window.location.host}/api/ws`;
    // let socket = new WebSocket(wsUrl);

    // socket.onmessage = (event) => {
    //     if (event.data === 'changed') {
    //         refreshBanner.style.display = 'block';
    //     }
    // };

    // socket.onclose = () => {
    //     console.log('WS connection closed. Reconnecting in 5s...');
    //     setTimeout(() => {
    //         socket = new WebSocket(wsUrl);
    //     }, 5000);
    // };

    // refreshBanner.onclick = () => {
    //     loadExpenses();
    //     refreshBanner.style.display = 'none';
    // };

    // Show modal for adding
    addBtn.onclick = () => {
        expenseForm.reset();
        document.getElementById('expense-id').value = '';
        modalTitle.innerText = 'Add Expense';
        document.getElementById('date').value = new Date().toISOString().split('T')[0];
        modal.style.display = 'block';
    };

    // Close modal
    closeBtn.onclick = () => modal.style.display = 'none';
    window.onclick = (event) => {
        if (event.target == modal) modal.style.display = 'none';
    };

    // Form submission
    expenseForm.onsubmit = async (e) => {
        e.preventDefault();
        const id = document.getElementById('expense-id').value;
        const data = {
            date: new Date(document.getElementById('date').value).toISOString(),
            description: document.getElementById('description').value,
            amount: parseFloat(document.getElementById('amount').value)
        };

        try {
            const method = id ? 'PUT' : 'POST';
            const url = id ? `${apiRoot}/expenses/${id}` : `${apiRoot}/expenses`;

            const response = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                const error = await response.text();
                throw new Error(error);
            }

            modal.style.display = 'none';
            loadExpenses();
        } catch (err) {
            alert('Error: ' + err.message);
        }
    };

    async function loadExpenses() {
        try {
            const response = await fetch(`${apiRoot}/expenses`);
            const expenses = await response.json();

            expenseList.innerHTML = '';
            if (!expenses) return;

            expenses.sort((a, b) => new Date(b.Date) - new Date(a.Date));

            expenses.forEach(e => {
                const row = document.createElement('tr');
                const date = new Date(e.Date).toLocaleDateString();
                row.innerHTML = `
                    <td>${date}</td>
                    <td>${e.Description}</td>
                    <td class="amount">Rs.${e.Amount.toFixed(2)}</td>
                    <td class="actions">
                        <button class="btn btn-small edit-btn" onclick="editExpense(${e.Id})">Edit</button>
                        <button class="btn btn-small delete-btn" onclick="deleteExpense(${e.Id})">Delete</button>
                    </td>
                `;
                expenseList.appendChild(row);
            });
        } catch (err) {
            console.error('Failed to load expenses', err);
        }
    }

    window.editExpense = async (id) => {
        try {
            const response = await fetch(`${apiRoot}/expenses/${id}`);
            const e = await response.json();

            document.getElementById('expense-id').value = e.Id;
            document.getElementById('date').value = e.Date.split('T')[0];
            document.getElementById('description').value = e.Description;
            document.getElementById('amount').value = e.Amount;

            modalTitle.innerText = 'Edit Expense';
            modal.style.display = 'block';
        } catch (err) {
            alert('Failed to fetch expense details');
        }
    };

    window.deleteExpense = async (id) => {
        if (!confirm('Are you sure you want to delete this expense?')) return;

        try {
            await fetch(`${apiRoot}/expenses/${id}`, { method: 'DELETE' });
            loadExpenses();
        } catch (err) {
            alert('Failed to delete expense');
        }
    };
});
