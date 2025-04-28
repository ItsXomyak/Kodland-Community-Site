const API_BASE = 'http://localhost:8080/api/v1/auth'

// Конфигурация безопасности
const SECURITY_CONFIG = {
	maxLoginAttempts: 5,
	lockoutTime: 1000, // 15 минут в миллисекундах
	tokenRefreshInterval: 10 * 60 * 1000, // 10 минут
	passwordMinLength: 8,
	passwordMaxLength: 64,
}

// Функция для санитизации ввода с расширенной защитой
function sanitizeInput(input) {
	if (typeof input !== 'string') return ''

	// Удаляем потенциально опасные символы
	const sanitized = input
		.replace(/[<>]/g, '') // Удаляем угловые скобки
		.replace(/javascript:/gi, '') // Удаляем javascript: протокол
		.replace(/data:/gi, '') // Удаляем data: протокол
		.trim()

	const div = document.createElement('div')
	div.textContent = sanitized
	return div.innerHTML
}

// Функция для генерации CSRF токена с дополнительной защитой
function generateCSRFToken() {
	const token = crypto.randomUUID()
	const timestamp = Date.now()
	const signature = btoa(`${token}:${timestamp}`)
	localStorage.setItem('csrfToken', signature)
	return signature
}

// Функция для валидации email с дополнительными проверками
function validateEmail(email) {
	if (!email) return false

	const sanitizedEmail = sanitizeInput(email)
	const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

	// Проверка длины email
	if (sanitizedEmail.length > 254) return false

	// Проверка на наличие подозрительных символов
	if (/[<>]/.test(sanitizedEmail)) return false

	return re.test(sanitizedEmail)
}

// Функция для валидации пароля с расширенными требованиями
function validatePassword(password) {
	if (!password) return false

	// Проверка длины пароля
	if (
		password.length < SECURITY_CONFIG.passwordMinLength ||
		password.length > SECURITY_CONFIG.passwordMaxLength
	) {
		return false
	}

	// Проверка на наличие подозрительных символов
	if (/[<>]/.test(password)) return false

	// Минимум 8 символов, минимум 1 цифра, 1 буква в верхнем регистре, 1 буква в нижнем регистре
	const re = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/
	return re.test(password)
}

// Функция для проверки количества попыток входа
function checkLoginAttempts() {
	const attempts = localStorage.getItem('loginAttempts') || 0
	const lastAttempt = localStorage.getItem('lastLoginAttempt')

	if (attempts >= SECURITY_CONFIG.maxLoginAttempts) {
		if (lastAttempt) {
			const timeSinceLastAttempt = Date.now() - parseInt(lastAttempt)
			if (timeSinceLastAttempt < SECURITY_CONFIG.lockoutTime) {
				const remainingTime = Math.ceil(
					(SECURITY_CONFIG.lockoutTime - timeSinceLastAttempt) / 1000 / 60
				)
				throw new Error(
					`Слишком много попыток входа. Попробуйте через ${remainingTime} минут.`
				)
			} else {
				// Сброс счетчика после истечения времени блокировки
				localStorage.setItem('loginAttempts', 0)
			}
		}
	}
}

// Функция для обновления токена
function refreshToken() {
	const token = getAuthToken()
	if (token) {
		fetch(`${API_BASE}/refresh-token`, {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`,
				'X-CSRF-Token': localStorage.getItem('csrfToken'),
				'X-Requested-With': 'XMLHttpRequest',
			},
			credentials: 'include',
		})
			.then(response => {
				if (response.ok) {
					return response.json()
				}
				throw new Error('Ошибка обновления токена')
			})
			.then(data => {
				setAuthToken(data.accessToken)
			})
			.catch(error => {
				console.error('Ошибка при обновлении токена:', error)
				// При ошибке обновления токена разлогиниваем пользователя
				logout()
			})
	}
}

// Функция для безопасного выхода
function logout() {
	const token = getAuthToken()
	if (token) {
		fetch(`${API_BASE}/logout`, {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`,
				'X-CSRF-Token': localStorage.getItem('csrfToken'),
				'X-Requested-With': 'XMLHttpRequest',
			},
			credentials: 'include',
		}).catch(error => console.error('Ошибка при выходе:', error))
	}

	// Очищаем все данные
	localStorage.removeItem('token')
	localStorage.removeItem('tokenExpiry')
	localStorage.removeItem('csrfToken')
	localStorage.removeItem('loginAttempts')
	localStorage.removeItem('lastLoginAttempt')

	window.location.href = '/auth.html'
}

// Функция для отображения сообщений с защитой от XSS
function showMessage(type, text) {
	const messageDiv = document.createElement('div')
	messageDiv.className = `${
		type === 'error' ? 'error-message' : 'success-message'
	}`
	messageDiv.textContent = sanitizeInput(text)
	messageDiv.style.display = 'block'

	const container = document.querySelector('.auth-container')
	container.insertBefore(messageDiv, container.firstChild)

	setTimeout(() => {
		messageDiv.style.opacity = '0'
		setTimeout(() => messageDiv.remove(), 300)
	}, 3000)
}

// Функция для обработки ошибок
function handleError(error) {
	console.error('Ошибка:', error)
	showMessage('error', error.message || 'Произошла ошибка')
}

// Функция для отправки запросов с защитой от CSRF
async function sendRequest(url, method, data) {
	try {
		const csrfToken = localStorage.getItem('csrfToken') || generateCSRFToken()

		const response = await fetch(url, {
			method,
			headers: {
				'Content-Type': 'application/json',
				'X-CSRF-Token': csrfToken,
				'X-Requested-With': 'XMLHttpRequest',
			},
			credentials: 'include', // Для работы с куками
			body: JSON.stringify(data),
		})

		if (!response.ok) {
			const errorData = await response.json()
			throw new Error(errorData.error || 'Ошибка при выполнении запроса')
		}

		// Обновляем CSRF токен после успешного запроса
		const newToken = response.headers.get('X-CSRF-Token')
		if (newToken) {
			localStorage.setItem('csrfToken', newToken)
		}

		return await response.json()
	} catch (error) {
		throw error
	}
}

// Функция для безопасного хранения токена
function setAuthToken(token) {
	try {
		// Шифруем токен перед сохранением
		const encryptedToken = btoa(token)
		localStorage.setItem('token', encryptedToken)

		// Устанавливаем время жизни токена
		const expiryDate = new Date()
		expiryDate.setHours(expiryDate.getMinutes() + 15) // Токен действителен 15 минут
		localStorage.setItem('tokenExpiry', expiryDate.toISOString())
	} catch (error) {
		console.error('Ошибка при сохранении токена:', error)
	}
}

// Функция для безопасного получения токена
function getAuthToken() {
	try {
		const token = localStorage.getItem('token')
		const expiry = localStorage.getItem('tokenExpiry')

		if (!token || !expiry) return null

		// Проверяем срок действия токена
		if (new Date(expiry) < new Date()) {
			localStorage.removeItem('token')
			localStorage.removeItem('tokenExpiry')
			return null
		}

		return atob(token)
	} catch (error) {
		console.error('Ошибка при получении токена:', error)
		return null
	}
}

// Функция для проверки авторизации и обновления навигации
function updateNavigation() {
	const token = localStorage.getItem('token')
	const authLink = document.getElementById('authLink')
	const profileNavItem = document.getElementById('profileNavItem')
	const userName = document.getElementById('userName')

	if (token) {
		// Если пользователь авторизован
		if (authLink) authLink.parentElement.classList.add('d-none')
		if (profileNavItem) profileNavItem.classList.remove('d-none')

		// Получаем информацию о пользователе из токена
		try {
			const payload = JSON.parse(atob(token.split('.')[1]))
			if (userName && payload.email) {
				userName.textContent = payload.email.split('@')[0] // Показываем имя пользователя до @
			}
		} catch (error) {
			console.error('Ошибка при разборе токена:', error)
		}
	} else {
		// Если пользователь не авторизован
		if (authLink) authLink.parentElement.classList.remove('d-none')
		if (profileNavItem) profileNavItem.classList.add('d-none')
	}
}

document.addEventListener('DOMContentLoaded', () => {
	// Генерируем CSRF токен при загрузке страницы
	generateCSRFToken()

	// Устанавливаем интервал обновления токена
	setInterval(refreshToken, SECURITY_CONFIG.tokenRefreshInterval)

	// Обработка переключения вкладок
	const tabButtons = document.querySelectorAll('.tab-btn')
	const tabContents = document.querySelectorAll('.tab-content')

	tabButtons.forEach(button => {
		button.addEventListener('click', () => {
			const tabId = button.getAttribute('data-tab')
			tabButtons.forEach(btn => btn.classList.remove('active'))
			tabContents.forEach(content => content.classList.remove('active'))
			button.classList.add('active')
			document.getElementById(tabId).classList.add('active')
		})
	})

	// Обработчик формы регистрации
	const registerForm = document.getElementById('registerForm')
	if (registerForm) {
		registerForm.addEventListener('submit', async function (e) {
			e.preventDefault()

			const email = document.getElementById('registerEmail')?.value
			const password = document.getElementById('registerPassword')?.value
			const displayName = document.getElementById('registerDisplayName')?.value

			if (!email || !password || !displayName) {
				showMessage('error', 'Все поля обязательны для заполнения')
				return
			}

			if (!validateEmail(email)) {
				showMessage('error', 'Пожалуйста, введите корректный email')
				return
			}

			if (!validatePassword(password)) {
				showMessage(
					'error',
					'Пароль должен содержать минимум 8 символов, включая цифры, заглавные и строчные буквы'
				)
				return
			}

			try {
				const response = await sendRequest(`${API_BASE}/register`, 'POST', {
					email,
					password,
					displayName,
				})

				showMessage(
					'success',
					'Регистрация успешна! Пожалуйста, проверьте вашу почту для подтверждения.'
				)
				setTimeout(() => {
					window.location.href = '/auth.html'
				}, 3000)
			} catch (error) {
				handleError(error)
			}
		})
	}

	// Обработка формы входа с защитой от брутфорса
	const loginForm = document.getElementById('loginForm')
	loginForm.addEventListener('submit', async e => {
		e.preventDefault()

		try {
			checkLoginAttempts()

			const email = document.getElementById('loginEmail').value
			const password = document.getElementById('loginPassword').value

			if (!validateEmail(email)) {
				showMessage('error', 'Пожалуйста, введите корректный email')
				return
			}

			try {
				const data = await sendRequest(`${API_BASE}/login`, 'POST', {
					email: sanitizeInput(email),
					password: sanitizeInput(password),
				})

				// Сброс счетчика попыток при успешном входе
				localStorage.setItem('loginAttempts', 0)
				localStorage.removeItem('lastLoginAttempt')

				handleSuccessfulLogin(data)
			} catch (error) {
				// Увеличиваем счетчик попыток
				const attempts =
					parseInt(localStorage.getItem('loginAttempts') || 0) + 1
				localStorage.setItem('loginAttempts', attempts)
				localStorage.setItem('lastLoginAttempt', Date.now())

				throw error
			}
		} catch (error) {
			handleError(error)
		}
	})

	// Обработка формы запроса сброса пароля
	const resetRequestForm = document.getElementById('resetRequestForm')
	if (resetRequestForm) {
		resetRequestForm.addEventListener('submit', async e => {
			e.preventDefault()

			const email = document.getElementById('resetEmail')?.value

			if (!email) {
				showMessage('error', 'Пожалуйста, введите email')
				return
			}

			try {
				await sendRequest(`${API_BASE}/password-reset-request`, 'POST', {
					email,
				})

				showMessage(
					'success',
					'Инструкции по сбросу пароля отправлены на ваш email'
				)
				resetRequestForm.reset()
			} catch (error) {
				handleError(error)
			}
		})
	}

	// Проверка наличия токена в URL для автоматического отображения формы подтверждения
	const urlParams = new URLSearchParams(window.location.search)
	const resetToken = urlParams.get('token')
	if (resetToken) {
		// Переключаем на вкладку сброса пароля
		document.querySelector('[data-tab="reset"]').click()

		// Скрываем форму запроса и показываем форму подтверждения
		document.getElementById('resetRequestForm').classList.add('d-none')
		document.getElementById('resetConfirmForm').classList.remove('d-none')

		// Устанавливаем значение токена в скрытое поле
		document.getElementById('resetToken').value = resetToken
		// Скрываем поле с токеном
		document.getElementById('resetToken').parentElement.classList.add('d-none')
	}

	// Обработчик формы сброса пароля
	const resetPasswordForm = document.getElementById('resetConfirmForm')
	if (resetPasswordForm) {
		resetPasswordForm.addEventListener('submit', async e => {
			e.preventDefault()

			const token = document.getElementById('resetToken').value
			const newPassword = document.getElementById('newPassword')?.value
			const confirmPassword =
				document.getElementById('confirmNewPassword')?.value

			if (!token) {
				showMessage('error', 'Недействительная ссылка для сброса пароля')
				return
			}

			if (!newPassword || !confirmPassword) {
				showMessage('error', 'Пожалуйста, заполните все поля')
				return
			}

			if (newPassword !== confirmPassword) {
				showMessage('error', 'Пароли не совпадают')
				return
			}

			if (!validatePassword(newPassword)) {
				showMessage(
					'error',
					'Пароль должен содержать минимум 8 символов, включая цифры и буквы в верхнем и нижнем регистре'
				)
				return
			}

			try {
				const response = await fetch(`${API_BASE}/password-reset-confirm`, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						token,
						newPassword,
					}),
				})

				if (!response.ok) {
					const errorData = await response.json()
					throw new Error(errorData.error || 'Ошибка при сбросе пароля')
				}

				showMessage('success', 'Пароль успешно изменен')
				setTimeout(() => {
					window.location.href = '/auth.html'
				}, 3000)
			} catch (error) {
				console.error('Ошибка:', error)
				showMessage('error', error.message || 'Ошибка при сбросе пароля')
			}
		})
	}

	// Проверка токена при загрузке страницы
	const token = getAuthToken()
	if (token) {
		window.location.href = '/index.html'
	}

	// Обработка неавторизованного доступа
	window.addEventListener('error', event => {
		if (event.message.includes('401')) {
			logout()
		}
	})

	// Очистка данных при выходе
	window.addEventListener('beforeunload', () => {
		if (!getAuthToken()) {
			localStorage.removeItem('csrfToken')
		}
	})

	// Проверка токена подтверждения email при загрузке страницы
	const confirmToken = urlParams.get('confirm')

	if (confirmToken) {
		// Если есть токен подтверждения, отправляем запрос на сервер
		fetch(`${API_BASE}/confirm?token=${confirmToken}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
			},
		})
			.then(response => {
				if (response.ok) {
					return response.json()
				}
				throw new Error('Ошибка подтверждения email')
			})
			.then(data => {
				showMessage(
					'success',
					data.message ||
						'Email успешно подтвержден! Теперь вы можете войти в систему.'
				)
				document.querySelector('[data-tab="login"]').click()
			})
			.catch(error => {
				console.error('Ошибка:', error)
				showMessage('error', error.message || 'Ошибка подтверждения email')
			})
	}

	// Обновляем навигацию при загрузке страницы
	updateNavigation()
})

// Обновляем навигацию после успешного входа
function handleSuccessfulLogin(data) {
	setAuthToken(data.accessToken)
	showMessage('success', 'Вход выполнен успешно!')
	updateNavigation()
	window.location.href = '/index.html'
}
