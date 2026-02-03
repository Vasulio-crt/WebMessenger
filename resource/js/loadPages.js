const searchInput = document.getElementById("searchInput");
const chatList = document.getElementById("chatList");
const searchButton = document.getElementById("searchButton");

function getPersonalChats() {
	const chats = localStorage.getItem("personalChats");
	return chats ? JSON.parse(chats) : [];
}

function savePersonalChats(chats) {
	localStorage.setItem("personalChats", JSON.stringify(chats));
}

function renderChats() {
	chatList.innerHTML = '';

	// Добавляем глобальный чат
	const globalChatLi = document.createElement("li");
	globalChatLi.innerHTML = `<a href="/globalChat">GLOBAL CHAT</a>`;
	if (window.location.pathname === '/globalChat') {
		globalChatLi.classList.add('active');
	}
	chatList.appendChild(globalChatLi);

	// Добавляем личные чаты
	const personalChats = getPersonalChats();
	personalChats.forEach(chatName => {
		const li = document.createElement("li");
		li.dataset.chatName = chatName;

		const chatLink = document.createElement("a");
		chatLink.href = `/chat/${chatName}`;
		chatLink.textContent = chatName;

		const deleteBtn = document.createElement("button");
		deleteBtn.innerHTML = `<img src="/resource/cross.svg" alt="Удалить чат" />`;
		deleteBtn.classList.add("delete-chat-btn");
		deleteBtn.onclick = (e) => {
			e.preventDefault();
			e.stopPropagation();
			removePersonalChat(chatName);
		};

		li.appendChild(chatLink);
		li.appendChild(deleteBtn);

		if (window.location.pathname === `/chat/${chatName}`) {
			li.classList.add('active');
		}

		chatList.appendChild(li);
	});
}

function addPersonalChat(chatName) {
	const personalChats = getPersonalChats();
	if (!personalChats.includes(chatName) && chatName !== userName) {
		personalChats.push(chatName);
		savePersonalChats(personalChats);
		renderChats();
	}
}

function removePersonalChat(chatName) {
	let personalChats = getPersonalChats();
	personalChats = personalChats.filter(name => name !== chatName);
	savePersonalChats(personalChats);
	renderChats();

	if (window.location.pathname === `/chat/${chatName}`) {
		window.location.href = '/globalChat';
	}
}

async function findAndAddUser(userName) {
	if (!userName.trim()) return;

	try {
		const response = await fetch(`/chat/find/${userName}`, {
			headers: {
				'X-Requested-With': 'XMLHttpRequest'
			}
		});
		const data = await response.json();

		if (data.found) {
			addPersonalChat(userName);
			searchInput.value = '';
		} else {
			alert("Пользователь не найден.");
		}
	} catch (error) {
		console.error("Ошибка при поиске пользователя:", error);
		alert("Произошла ошибка при поиске.");
	}
}

searchButton.addEventListener("click", function() {
	findAndAddUser(searchInput.value.trim());
});

document.addEventListener("DOMContentLoaded", function() {
	renderChats();
});
