const userNameElement = document.getElementById("userName");
if (userNameElement) {
	userNameElement.innerText = userName;
}

const Pathname = window.location.pathname;
const HistoryPath = Pathname + "/history";
const socket = new WebSocket(`ws://${window.location.host + Pathname}/ws`);
const messages = document.getElementById("messages");
const messageInput = document.getElementById("messageInput");
const sendMessageButton = document.getElementById("sendMessageButton");

function displayMessage(data) {
	const message = document.createElement("li");
	try {
		const parsedData = JSON.parse(data);
		if (parsedData.from === userName) {
			message.textContent = parsedData.text;
			message.classList.add("sent-mes");
		} else if(parsedData.from === "server") {
			message.textContent = parsedData.text;
			message.classList.add("server-mes");
		} else {
			message.textContent = `${parsedData.from}: ${parsedData.text}`;
			message.classList.add("received-mes");
		}
	} catch (e) {
		message.textContent = data;
		message.classList.add("received-mes");
	}
	messages.appendChild(message);
	messages.scrollTop = messages.scrollHeight;
}

socket.onopen = function(event) {
	console.log("WebSocket connection established.");
};

function showNotification(fromUser) {
	const notificationContainer = document.getElementById('notification-container');
	const notification = document.createElement('div');
	notification.className = 'notification';
	notification.textContent = `Новое сообщение от ${fromUser}`;

	notification.onclick = () => {
		window.location.href = `/chat/${fromUser}`;
	};

	notificationContainer.appendChild(notification);

	setTimeout(() => {
		notification.style.opacity = '0';
		setTimeout(() => notification.remove(), 500);
	}, 5000);
}

socket.onmessage = function(event) {
	try {
		const parsedData = JSON.parse(event.data);
		if (parsedData.from !== userName && Pathname !== `/chat/${parsedData.from}`) {
			showNotification(parsedData.from);
			return;
		}
	} catch (e) {
		console.error(e);
	}
	displayMessage(event.data);
};

socket.onclose = function(event) {
	console.log("WebSocket connection closed.");
	const message = document.createElement("li");
	message.textContent = "Connection closed.";
	message.style.textAlign = "center";
	message.style.color = "#888";
	messages.appendChild(message);
};

socket.onerror = function(error) {
	console.error("WebSocket Error: ", error);
};

// --- sendMessage ---
function sendMessage() {
	const messageText = messageInput.value.trim();
	if (messageText === "") {
		return;
	}

	let messagePayload;
	if (Pathname === "/globalChat") {
		messagePayload = JSON.stringify({
			from: userName,
			text: messageText
		});
	} else {
		messagePayload = JSON.stringify({
			from: userName,
			to: Pathname.slice(6),
			text: messageText
		});
		displayMessage(messagePayload);	
	}
	socket.send(messagePayload);
	messageInput.value = "";
}

messageInput.addEventListener("keydown", function(event) {
	if (event.key === "Enter") {
		event.preventDefault();
		sendMessage();
	}
});
sendMessageButton.addEventListener("click", sendMessage);

// --- logout ---
document.getElementById('logoutButton').addEventListener('click', async function() {
	try {
		const response = await fetch('/logout', {
			method: 'POST',
		});
		localStorage.removeItem('userName');
		if (response.ok) {
			window.location.href = '/login';
		}
	} catch (error) {
		console.error('Error sending logout request:', error);
	}
});

// --- burgerMenu ---
const burgerMenu = document.getElementById('burger-menu');
const sidebar = document.querySelector('.sidebar');

burgerMenu.addEventListener('click', function() {
	sidebar.classList.add('open');
	burgerMenu.style.display = 'none';
});

document.querySelector('.main-content').addEventListener('click', function(e) {
	if (sidebar.classList.contains('open') && !e.target.closest('.sidebar') && e.target !== burgerMenu) {
		sidebar.classList.remove('open')
		setTimeout(() => burgerMenu.style.display = 'block', 300);
	};
});

// --- loadHistory ---
function loadHistory() {
	fetch(HistoryPath)
		.then(response => {
			if (!response.ok) {
				throw new Error("Network response was not ok");
			}
			return response.json();
		})
		.then(history => {
			if (history === null) {
				console.log("History is empty.");
			} else {
				messages.innerHTML = ''; 
				history.forEach((obj) => displayMessage(JSON.stringify(obj)));
				console.log("History loaded.");
			}
		})
		.catch(error => {
			console.error("Error loading history:", error);
			alert("Не удалось загрузить историю сообщений.");
		});
};

window.onload = function () {
	loadHistory();
};
