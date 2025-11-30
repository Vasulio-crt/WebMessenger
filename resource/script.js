let userId = localStorage.getItem("userId");
if (!userId) {
	userId = "user-" + Math.random().toString(16).slice(2);
	localStorage.setItem("userId", userId);
}

document.getElementById("userName").innerText = userId;

const IP = "10.136.119.49";

const socket = new WebSocket(`ws://${IP}:8080/ws`);
const messages = document.getElementById("messages");
const messageInput = document.getElementById("messageInput");

function displayMessage(data) {
	const message = document.createElement("li");
	try {
		const parsedData = JSON.parse(data);
		message.textContent = parsedData.text;

		if (parsedData.senderId === userId) {
			message.classList.add("sent-mes");
		} else if(parsedData.senderId === "server") {
			message.classList.add("server-mes");
		} else {
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

socket.onmessage = function(event) {
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

messageInput.addEventListener("keydown", function(event) {
	if (event.key === "Enter") {
		event.preventDefault();
		const messageText = messageInput.value.trim();
		if (messageText !== "") {
			const messagePayload = JSON.stringify({
				senderId: userId,
				text: messageText
			});
			socket.send(messagePayload);
			messageInput.value = "";
		}
	}
});

function loadHistory() {
	fetch("/history")
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
// window.addEventListener("beforeunload", function (e) {
// 	e.preventDefault();
// 	e.returnValue = 'Сообщения будут потеряны.';
// });
