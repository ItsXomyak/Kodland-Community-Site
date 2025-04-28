// Константы
const API_BASE = 'http://localhost:8080/api/v1'
const PROFILE_ENDPOINT = `${API_BASE}/profile`

// DOM элементы
const profileAvatar = document.getElementById('profileAvatar')
const profileName = document.getElementById('profileName')
const profileEmail = document.getElementById('profileEmail')
const editProfileModal = new bootstrap.Modal(
	document.getElementById('editProfileModal')
)
const editProfileForm = document.getElementById('editProfileForm')
const editFirstName = document.getElementById('editFirstName')
const editLastName = document.getElementById('editLastName')
const editBio = document.getElementById('editBio')
const passwordForm = document.getElementById('passwordForm')
const avatarInput = document.getElementById('avatarInput')
const avatarPreview = document.getElementById('avatarPreview')
const avatarEditBtn = document.getElementById('avatarEditBtn')
const statsContainer = document.getElementById('statsContainer')
const activityList = document.getElementById('activityList')
const achievementsGrid = document.getElementById('achievementsGrid')
const tablesList = document.getElementById('tablesList')

// Загрузка данных профиля
async function loadProfileData() {
	try {
		const token = localStorage.getItem('token')
		if (!token) {
			window.location.href = '/auth.html'
			return
		}

		const response = await fetch(PROFILE_ENDPOINT, {
			headers: {
				Authorization: `Bearer ${token}`,
			},
		})

		if (!response.ok) {
			if (response.status === 401) {
				localStorage.removeItem('token')
				window.location.href = '/auth.html'
				return
			}
			throw new Error('Ошибка загрузки данных профиля')
		}

		const data = await response.json()
		updateProfileUI(data)
	} catch (error) {
		console.error('Ошибка:', error)
		showError('Не удалось загрузить данные профиля')
	}
}

// Обновление UI профиля
function updateProfileUI(data) {
	// Обновление основной информации
	profileName.textContent = `${data.firstName} ${data.lastName}`
	profileEmail.textContent = data.email

	// Обновление аватара
	if (data.avatar) {
		profileAvatar.src = data.avatar
	}

	// Заполнение формы редактирования
	editFirstName.value = data.firstName
	editLastName.value = data.lastName
	editBio.value = data.bio || ''

	// Обновление статистики
	updateStats(data.stats)

	// Обновление активности
	updateActivity(data.activity)

	// Обновление достижений
	updateAchievements(data.achievements)

	// Обновление таблиц
	updateTables(data.tables)
}

// Обновление статистики
function updateStats(stats) {
	const statsHTML = `
        <div class="stat-item">
            <i class="fas fa-medal"></i>
            <span>${stats.medals || 0}</span>
            <small>Медали</small>
        </div>
        <div class="stat-item">
            <i class="fas fa-trophy"></i>
            <span>${stats.achievements || 0}</span>
            <small>Достижения</small>
        </div>
        <div class="stat-item">
            <i class="fas fa-table"></i>
            <span>${stats.tables || 0}</span>
            <small>Таблицы</small>
        </div>
    `
	statsContainer.innerHTML = statsHTML
}

// Обновление активности
function updateActivity(activity) {
	if (!activity || activity.length === 0) {
		activityList.innerHTML = '<p class="text-muted">Нет активности</p>'
		return
	}

	const activityHTML = activity
		.map(
			item => `
        <div class="activity-item">
            <div class="activity-icon">
                <i class="fas ${getActivityIcon(item.type)}"></i>
            </div>
            <div class="activity-content">
                <p class="mb-1">${item.description}</p>
                <small class="activity-date">${formatDate(item.date)}</small>
            </div>
        </div>
    `
		)
		.join('')
	activityList.innerHTML = activityHTML
}

// Обновление достижений
function updateAchievements(achievements) {
	if (!achievements || achievements.length === 0) {
		achievementsGrid.innerHTML = '<p class="text-muted">Нет достижений</p>'
		return
	}

	const achievementsHTML = achievements
		.map(
			achievement => `
        <div class="achievement-card">
            <div class="achievement-icon">
                <i class="fas ${achievement.icon}"></i>
            </div>
            <h5 class="achievement-title">${achievement.title}</h5>
            <p class="achievement-description">${achievement.description}</p>
        </div>
    `
		)
		.join('')
	achievementsGrid.innerHTML = achievementsHTML
}

// Обновление таблиц
function updateTables(tables) {
	if (!tables || tables.length === 0) {
		tablesList.innerHTML = '<p class="text-muted">Нет таблиц</p>'
		return
	}

	const tablesHTML = tables
		.map(
			table => `
        <div class="table-card">
            <div class="table-header">
                <h5 class="table-title">${table.title}</h5>
                <div class="table-meta">
                    <span><i class="fas fa-calendar"></i> ${formatDate(
											table.createdAt
										)}</span>
                    <span><i class="fas fa-eye"></i> ${table.views}</span>
                </div>
            </div>
            <p class="mb-2">${table.description}</p>
            <div class="table-actions">
                <button class="btn btn-sm btn-primary" onclick="viewTable('${
									table.id
								}')">
                    <i class="fas fa-eye"></i> Просмотр
                </button>
                <button class="btn btn-sm btn-secondary" onclick="editTable('${
									table.id
								}')">
                    <i class="fas fa-edit"></i> Редактировать
                </button>
            </div>
        </div>
    `
		)
		.join('')
	tablesList.innerHTML = tablesHTML
}

// Функция открытия модального окна редактирования
function editProfile() {
	editProfileModal.show()
}

// Функция сохранения профиля
async function saveProfile() {
	const data = {
		firstName: editFirstName.value,
		lastName: editLastName.value,
		bio: editBio.value,
	}

	try {
		const token = localStorage.getItem('token')
		if (!token) {
			window.location.href = '/auth.html'
			return
		}

		const response = await fetch(PROFILE_ENDPOINT, {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`,
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(data),
		})

		if (!response.ok) {
			if (response.status === 401) {
				localStorage.removeItem('token')
				window.location.href = '/auth.html'
				return
			}
			throw new Error('Ошибка обновления профиля')
		}

		showSuccess('Профиль успешно обновлен')
		editProfileModal.hide()
		loadProfileData()
	} catch (error) {
		console.error('Ошибка:', error)
		showError('Не удалось обновить профиль')
	}
}

// Обработка формы смены пароля
passwordForm.addEventListener('submit', async e => {
	e.preventDefault()

	const formData = new FormData(passwordForm)
	const data = {
		currentPassword: formData.get('currentPassword'),
		newPassword: formData.get('newPassword'),
	}

	try {
		const response = await fetch(`${PROFILE_ENDPOINT}/password`, {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${localStorage.getItem('token')}`,
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(data),
		})

		if (!response.ok) {
			throw new Error('Ошибка смены пароля')
		}

		showSuccess('Пароль успешно изменен')
		passwordForm.reset()
	} catch (error) {
		console.error('Ошибка:', error)
		showError('Не удалось изменить пароль')
	}
})

// Обработка загрузки аватара
avatarInput.addEventListener('change', async e => {
	const file = e.target.files[0]
	if (!file) return

	const formData = new FormData()
	formData.append('avatar', file)

	try {
		const response = await fetch(`${PROFILE_ENDPOINT}/avatar`, {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${localStorage.getItem('token')}`,
			},
			body: formData,
		})

		if (!response.ok) {
			throw new Error('Ошибка загрузки аватара')
		}

		const data = await response.json()
		avatarPreview.src = data.avatarUrl
		showSuccess('Аватар успешно обновлен')
	} catch (error) {
		console.error('Ошибка:', error)
		showError('Не удалось загрузить аватар')
	}
})

// Вспомогательные функции
function getActivityIcon(type) {
	const icons = {
		table_created: 'fa-table',
		achievement_unlocked: 'fa-trophy',
		medal_earned: 'fa-medal',
		profile_updated: 'fa-user-edit',
	}
	return icons[type] || 'fa-info-circle'
}

function formatDate(dateString) {
	const date = new Date(dateString)
	return date.toLocaleDateString('ru-RU', {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
	})
}

function showSuccess(message) {
	// Временное решение - можно заменить на более красивые уведомления
	alert(message)
}

function showError(message) {
	// Временное решение - можно заменить на более красивые уведомления
	alert(message)
}

// Инициализация
document.addEventListener('DOMContentLoaded', () => {
	loadProfileData()
})
