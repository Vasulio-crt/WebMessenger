const HOST = 'http://' + window.location.host;

/**
 * Настраивает обработчик отправки формы для аутентификации (регистрация или вход).
 * @param {string} formId ID HTML-формы.
 * @param {string} endpoint URL-адрес API для отправки данных (например, '/register' или '/login').
 * @param {string} successRedirectPath Путь для перенаправления после успешной аутентификации.
*/
function setupAuthForm(formId, endpoint, successRedirectPath) {
	document.getElementById(formId).addEventListener('submit', async function(event) {
		event.preventDefault();
		const userNameInput = document.getElementById('user_name');
		const passwordInput = document.getElementById('password');
		const messageDiv = document.getElementById('message');
		const userNameVal = userNameInput.value.trim();
		const passwordVal = passwordInput.value.trim();
		const userNameRegex = /^[a-zA-Z0-9_]{3,32}$/;
		
		if (!userNameVal || !passwordVal) {
			if (messageDiv) {
				messageDiv.textContent = 'Заполните все поля.';
			}
			return;
		}
		if (userNameVal.length > 32) {
			messageDiv.textContent = 'Имя пользователя должно быть от 3 до 32 символов.';
			return;
		}
		if (!userNameRegex.test(userNameVal)) {
			messageDiv.textContent = 'Имя пользователя может содержать только латинские буквы (a-z, A-Z), цифры (0-9) и знак подчеркивания (_).';
			return;
		}
		if (passwordVal.length > 32) {
			messageDiv.textContent = 'Пароль не должен превышать 32 символа.';
			return;
		}
		try {
			const formData = {
				userName: userNameVal,
				password: passwordVal
			};
			const response = await fetch(HOST + endpoint, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json;charset=utf-8'
				},
				body: JSON.stringify(formData)
			});
			if (response.ok) {
				localStorage.setItem('userName', userNameVal);
				window.location.href = HOST + successRedirectPath;
			} else {
				const errorText = await response.text();
				messageDiv.textContent = `Ошибка: ${errorText}`;
			}
		} catch (error) {
			console.error('Ошибка:', error);
			messageDiv.textContent = 'Ошибка соединения с сервером.';
		}
	});
}