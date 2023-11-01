document.addEventListener('DOMContentLoaded', function () {
    const messagesDiv = document.getElementById('messages');
    const messageInput = document.getElementById('message');
    const usernameInput = document.getElementById('username');

    function showError(error) {
        console.error(error);
        alert("Sorry, something went wrong. The API Service is not responding.");
    }

    function fetchMessages() {
        fetch(`${apiUrl}/messages`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok.');
                }
                return response.json();
            })
            .then(data => {
                messagesDiv.innerHTML = data.map(msg => `<div><strong>${msg.AuthorUsername}</strong>: ${msg.Text}</div>`).join('');
            })
            .catch(
                () => {
                    // do nothing on the poller
                }
            );
    }

    document.getElementById('send').addEventListener('click', () => {
        const message = messageInput.value;
        const username = usernameInput.value;

        if (!message || !username) {
            alert("Username and message are required!");
            return;
        }

        fetch(`${apiUrl}/messages`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ AuthorUsername: username, Text: message })
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok.');
                }
                return response.json();
            })
            .then(() => {
                messageInput.value = '';
                fetchMessages();
            })
            .catch(showError);
    });

    // Fetch messages initially and set an interval for refreshing
    fetchMessages();
    setInterval(fetchMessages, 5000); // Fetch every 5 seconds
});
