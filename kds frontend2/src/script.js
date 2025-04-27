// Глобальные переменные
let members = []
let currentPage = 1
const membersPerPage = 12
let isLoading = false
const API_URL = 'http://localhost:8080/api'

// Глобальная функция для обработки лайков
window.toggleLike = async function (event, userId) {
	if (event) {
		event.stopPropagation()
	}

	if (!userId) {
		console.error('ID пользователя не указан')
		return
	}

	try {
		const response = await fetch(`${API_URL}/users/${userId}/like`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
		})

		if (!response.ok) {
			throw new Error('Ошибка при обработке лайка')
		}

		const data = await response.json()

		// Обновляем количество лайков в модальном окне
		const modalLikes = document.getElementById('modalLikes')
		if (modalLikes) {
			modalLikes.textContent = data.likes
		}

		// Обновляем количество лайков в карточке
		const cardLikes = document.querySelector(
			`[data-user-id="${userId}"] .likes-count`
		)
		if (cardLikes) {
			cardLikes.textContent = data.likes
		}

		// Обновляем состояние кнопки лайка в модальном окне
		const modalLikeBtn = document.querySelector(
			'#modalLikesContainer .like-button'
		)
		if (modalLikeBtn) {
			modalLikeBtn.classList.toggle('text-red-500')
			modalLikeBtn.classList.toggle('text-gray-500')
		}

		// Обновляем состояние кнопки лайка в карточке
		const cardLikeBtn = document.querySelector(`[data-user-id="${userId}"]`)
		if (cardLikeBtn) {
			cardLikeBtn.classList.toggle('text-red-500')
			cardLikeBtn.classList.toggle('text-gray-500')
		}
	} catch (error) {
		console.error('Ошибка:', error)
	}
}

// Функция для загрузки данных с сервера
async function fetchMembers(page = 1) {
	if (isLoading) return
	isLoading = true
	showLoadingIndicator()

	try {
		const searchQuery = document
			.getElementById('searchInput')
			.value.toLowerCase()
		const roleFilter = document.getElementById('roleFilter').value
		const ageFilter = document.getElementById('ageFilter').value
		const countryFilter = document.getElementById('countryFilter').value
		const achievementFilter = document.getElementById('achievementFilter').value
		const sortBy = document.getElementById('sortSelect').value

		// Формируем параметры запроса
		const params = new URLSearchParams({
			page: page,
			limit: membersPerPage,
			search: searchQuery,
			role: roleFilter,
			age: ageFilter,
			country: countryFilter,
			achievement: achievementFilter,
			sort: sortBy,
		})

		const response = await fetch(`${API_URL}/users?${params}`)
		if (!response.ok) {
			throw new Error('Ошибка загрузки данных')
		}

		const data = await response.json()
		members = data.users
		const totalPages = Math.ceil(data.total / membersPerPage)

		// Отображаем участников
		displayMembers(members)

		// Настраиваем пагинацию
		setupPagination(totalPages)
	} catch (error) {
		showError(error.message)
	} finally {
		isLoading = false
		hideLoadingIndicator()
	}
}

// Функция для отображения индикатора загрузки
function showLoadingIndicator() {
	const grid = document.getElementById('membersGrid')
	grid.innerHTML =
		'<div class="col-span-full text-center py-8">Загрузка...</div>'
}

// Функция для скрытия индикатора загрузки
function hideLoadingIndicator() {
	// Индикатор скроется автоматически при отображении данных
}

// Функция для отображения ошибки
function showError(message) {
	const grid = document.getElementById('membersGrid')
	grid.innerHTML = `<div class="col-span-full text-center py-8 text-red-500">${message}</div>`
}

// Функция для настройки пагинации
function setupPagination(totalPages) {
	const paginationContainer = document.createElement('div')
	paginationContainer.className = 'flex justify-center mt-8 space-x-2'

	for (let i = 1; i <= totalPages; i++) {
		const button = document.createElement('button')
		button.className = `px-4 py-2 rounded ${
			i === currentPage
				? 'bg-kodland text-white'
				: 'bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-white'
		}`
		button.textContent = i
		button.onclick = () => {
			currentPage = i
			fetchMembers(i)
		}
		paginationContainer.appendChild(button)
	}

	const main = document.querySelector('main')
	const existingPagination = document.querySelector('.pagination-container')
	if (existingPagination) {
		existingPagination.remove()
	}
	paginationContainer.classList.add('pagination-container')
	main.appendChild(paginationContainer)
}

// Функция для создания карточки участника
function createMemberCard(member) {
	const card = document.createElement('div')
	card.className =
		'bg-white dark:bg-gray-800 rounded-lg shadow-lg p-4 cursor-pointer hover:shadow-xl transition-all duration-200 hover:scale-[1.02]'
	card.innerHTML = `
		<div class="flex items-center space-x-4">
			<div class="relative">
				<img src="${member.avatar_url}" alt="${member.nickname}'s avatar" 
					class="w-16 h-16 rounded-full object-cover border-2 border-kodland">
				${
					member.roles && member.roles.includes('Модератор')
						? `
					<div class="absolute -bottom-1 -right-1 bg-kodland text-white rounded-full p-1">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
						</svg>
					</div>
				`
						: ''
				}
			</div>
			<div class="flex-1 min-w-0">
				<h2 class="text-lg font-bold text-gray-800 dark:text-white truncate">${
					member.nickname
				}</h2>
				<div class="flex items-center space-x-2 mt-1">
					${
						member.medals &&
						(member.medals.gold > 0 ||
							member.medals.silver > 0 ||
							member.medals.bronze > 0)
							? `
						<div class="flex space-x-1">
          ${
						member.medals.gold > 0
							? `<span class="text-yellow-500">🥇${member.medals.gold}</span>`
							: ''
					}
          ${
						member.medals.silver > 0
							? `<span class="text-gray-400">🥈${member.medals.silver}</span>`
							: ''
					}
							${
								member.medals.bronze > 0
									? `<span class="text-orange-400">🥉${member.medals.bronze}</span>`
									: ''
							}
						</div>
					`
							: ''
					}
          ${
						member.achievements && member.achievements.length > 0
							? `
						<span class="text-kodland">🏆${member.achievements.length}</span>
					`
							: ''
					}
        </div>
      </div>
			<div class="flex items-center space-x-1 text-gray-500 dark:text-gray-400">
				<button class="like-button p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200 ${
					member.liked ? 'text-red-500' : 'text-gray-500'
				}" 
					data-user-id="${member.id}" 
					onclick="toggleLike(event, '${member.id}')">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path>
					</svg>
				</button>
				<span class="likes-count text-sm">${member.likes || 0}</span>
      </div>
    </div>
  `
	card.addEventListener('click', e => {
		// Предотвращаем открытие модального окна при клике на кнопку лайка
		if (!e.target.closest('.like-button')) {
			showUserModal(member.id)
		}
	})
	return card
}

// Функция для отображения всех карточек
function displayMembers(membersToShow = members) {
	const grid = document.getElementById('membersGrid')
	grid.innerHTML = '' // Очищаем сетку
	membersToShow.forEach(member => {
		const card = createMemberCard(member)
		grid.appendChild(card)
	})
}

// Функция для извлечения возраста из ролей
function extractAgeFromRoles(roles) {
	const agePattern = /(\d+)\+?\s*лет/
	for (const role of roles) {
		const match = role.match(agePattern)
		if (match) {
			return parseInt(match[1])
		}
	}
	return null
}

// Функция для извлечения страны из ролей
function extractCountryFromRoles(roles) {
	const countries = ['Германия', 'Франция', 'Россия']
	for (const role of roles) {
		if (countries.includes(role)) {
			return role
		}
	}
	return null
}

// Функция для фильтрации участников
function filterMembers() {
	const searchQuery = document.getElementById('searchInput').value.toLowerCase()
	const roleFilter = document.getElementById('roleFilter').value
	const ageFilter = document.getElementById('ageFilter').value
	const countryFilter = document.getElementById('countryFilter').value
	const achievementFilter = document.getElementById('achievementFilter').value
	const sortBy = document.getElementById('sortSelect').value

	let filteredMembers = members.filter(member => {
		// Поиск по никнейму
		if (searchQuery && !member.nickname.toLowerCase().includes(searchQuery)) {
			return false
		}

		// Фильтр по ролям
		if (roleFilter) {
			const rolesWithoutMeta = member.roles.filter(
				role =>
					!role.match(/\d+\+?\s*лет/) &&
					!['Германия', 'Франция', 'Россия'].includes(role)
			)
			if (!rolesWithoutMeta.includes(roleFilter)) {
				return false
			}
		}

		// Фильтр по возрасту
		if (ageFilter) {
			const memberAge = extractAgeFromRoles(member.roles)
			if (memberAge === null) return false
			if (ageFilter === '16+ лет' && memberAge < 16) return false
			if (ageFilter === '17 лет' && memberAge !== 17) return false
			if (ageFilter === '18+ лет' && memberAge < 18) return false
		}

		// Фильтр по стране
		if (countryFilter) {
			const memberCountry = extractCountryFromRoles(member.roles)
			if (memberCountry !== countryFilter) {
				return false
			}
		}

		// Фильтр по достижениям
		if (achievementFilter && !member.achievements.includes(achievementFilter)) {
			return false
		}

		return true
	})

	// Сортировка
	filteredMembers = sortMembers(filteredMembers, sortBy)

	// Отображение отфильтрованных участников
	displayMembers(filteredMembers)
}

// Функция для сортировки участников
function sortMembers(members, sortBy) {
	return [...members].sort((a, b) => {
		switch (sortBy) {
			case 'recent':
				return new Date(b.created_at) - new Date(a.created_at)
			case 'nickname':
				return a.nickname.localeCompare(b.nickname)
			case 'time':
				return new Date(b.joined_at) - new Date(a.joined_at)
			case 'achievements':
				return b.achievements.length - a.achievements.length
			case 'likes':
				return b.likes - a.likes
			case 'medals':
				const aTotalMedals =
					a.medals.gold * 3 + a.medals.silver * 2 + a.medals.bronze
				const bTotalMedals =
					b.medals.gold * 3 + b.medals.silver * 2 + b.medals.bronze
				return bTotalMedals - aTotalMedals
			default:
				return 0
		}
	})
}

// Функция для показа модального окна с полной информацией
function showUserModal(userId) {
	const member = members.find(m => m.id === userId)
	if (!member) return

	const modal = document.getElementById('userModal')
	if (!modal) {
		console.error('Модальное окно не найдено')
		return
	}

	// Заполняем базовую информацию
	const modalAvatar = modal.querySelector('#modalAvatar')
	const modalNickname = modal.querySelector('#modalNickname')
	const modalJoinDate = modal.querySelector('#modalJoinDate')
	const modalLikes = modal.querySelector('#modalLikes')
	const modalCreatedAt = modal.querySelector('#modalCreatedAt')

	if (modalAvatar) {
		modalAvatar.src = member.avatar_url
		modalAvatar.alt = `${member.nickname}'s avatar`
	}
	if (modalNickname) modalNickname.textContent = member.nickname
	if (modalJoinDate)
		modalJoinDate.textContent = `В сообществе с ${new Date(
			member.joined_at
		).toLocaleDateString('ru-RU')}`
	if (modalLikes) modalLikes.textContent = member.likes || '0'
	if (modalCreatedAt)
		modalCreatedAt.textContent = new Date(member.created_at).toLocaleDateString(
			'ru-RU'
		)

	// Добавляем отображение лайков в модальное окно
	const modalLikesContainer = modal.querySelector('#modalLikesContainer')
	if (modalLikesContainer) {
		modalLikesContainer.innerHTML = `
			<div class="flex items-center space-x-2">
				<span class="text-gray-500 dark:text-gray-400">❤️</span>
				<span class="likes-count text-sm">${member.likes || 0}</span>
			</div>
		`
	}

	// Создаем контейнер для дополнительной информации
	const modalContent = document.createElement('div')
	modalContent.className = 'mt-6 space-y-6'

	// Добавляем ссылки
	if (member.links && member.links.length > 0) {
		const linksSection = document.createElement('div')
		linksSection.innerHTML = `
			<h3 class="text-lg font-semibold text-gray-800 dark:text-white mb-2">Ссылки</h3>
			<div class="flex flex-wrap gap-2">
				${member.links
					.map(
						link => `
					<a href="${link}" target="_blank" class="text-kodland hover:underline flex items-center space-x-1">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path>
						</svg>
						<span>${new URL(link).hostname}</span>
					</a>
				`
					)
					.join('')}
      </div>
		`
		modalContent.appendChild(linksSection)
	}

	// Добавляем достижения
	if (member.achievements && member.achievements.length > 0) {
		const achievementsSection = document.createElement('div')
		achievementsSection.innerHTML = `
			<h3 class="text-lg font-semibold text-gray-800 dark:text-white mb-2">Достижения</h3>
			<div class="flex flex-wrap gap-2">
				${member.achievements
					.map(
						ach => `
					<span class="px-3 py-1 bg-kodland text-white rounded-full text-sm flex items-center">
						<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
						${ach}
					</span>
				`
					)
					.join('')}
			</div>
		`
		modalContent.appendChild(achievementsSection)
	}

	// Добавляем медали
	if (
		member.medals &&
		(member.medals.gold > 0 ||
			member.medals.silver > 0 ||
			member.medals.bronze > 0)
	) {
		const medalsSection = document.createElement('div')
		medalsSection.innerHTML = `
			<h3 class="text-lg font-semibold text-gray-800 dark:text-white mb-2">Медали</h3>
			<div class="flex flex-wrap gap-2">
				${
					member.medals.gold > 0
						? `
					<span class="px-3 py-1 bg-yellow-500 text-white rounded-full text-sm flex items-center">
						<span class="mr-1">${member.medals.gold}</span>
						<span>🥇</span>
					</span>
				`
						: ''
				}
				${
					member.medals.silver > 0
						? `
					<span class="px-3 py-1 bg-gray-400 text-white rounded-full text-sm flex items-center">
						<span class="mr-1">${member.medals.silver}</span>
						<span>🥈</span>
					</span>
				`
						: ''
				}
				${
					member.medals.bronze > 0
						? `
					<span class="px-3 py-1 bg-orange-400 text-white rounded-full text-sm flex items-center">
						<span class="mr-1">${member.medals.bronze}</span>
						<span>🥉</span>
					</span>
				`
						: ''
				}
			</div>
		`
		modalContent.appendChild(medalsSection)
	}

	// Добавляем роли
	if (member.roles && member.roles.length > 0) {
		const rolesSection = document.createElement('div')
		rolesSection.innerHTML = `
			<h3 class="text-lg font-semibold text-gray-800 dark:text-white mb-2">Роли</h3>
			<div class="flex flex-wrap gap-2">
				${member.roles
					.map(
						role => `
					<span class="px-3 py-1 bg-gray-200 dark:bg-gray-700 rounded-full text-sm text-gray-700 dark:text-gray-300 flex items-center">
						<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
						</svg>
        ${role}
      </span>
    `
					)
					.join('')}
			</div>
		`
		modalContent.appendChild(rolesSection)
	}

	// Очищаем предыдущий контент и добавляем новый
	const modalBody = modal.querySelector('.bg-white')
	const existingContent = modalBody.querySelector('.mt-6')
	if (existingContent) {
		existingContent.remove()
	}
	modalBody.appendChild(modalContent)

	// Показываем модальное окно
	modal.classList.remove('hidden')
	modal.classList.add('flex')
	setTimeout(() => {
		modal.querySelector('.bg-white').classList.remove('scale-95', 'opacity-0')
	}, 10)
}

// Функция для закрытия модального окна
function closeModal() {
	const modal = document.getElementById('userModal')
	if (modal) {
		modal.classList.remove('flex')
		modal.classList.add('hidden')
		const modalContent = modal.querySelector('.bg-white')
		if (modalContent) {
			modalContent.classList.remove('scale-100', 'opacity-100')
			modalContent.classList.add('scale-95', 'opacity-0')
		}
	}
}

// Функция для сброса всех фильтров
function resetFilters() {
	document.getElementById('searchInput').value = ''
	document.getElementById('roleFilter').value = ''
	document.getElementById('ageFilter').value = ''
	document.getElementById('countryFilter').value = ''
	document.getElementById('achievementFilter').value = ''
	document.getElementById('sortSelect').value = 'recent'

	// Применяем сброс фильтров
	filterMembers()

	// Добавляем анимацию сброса
	const resetButton = document.getElementById('resetFilters')
	resetButton.classList.add('animate-pulse')
	setTimeout(() => {
		resetButton.classList.remove('animate-pulse')
	}, 1000)
}

// Инициализация обработчиков событий
function initializeEventListeners() {
	const searchInput = document.getElementById('searchInput')
	const roleFilter = document.getElementById('roleFilter')
	const ageFilter = document.getElementById('ageFilter')
	const countryFilter = document.getElementById('countryFilter')
	const achievementFilter = document.getElementById('achievementFilter')
	const sortSelect = document.getElementById('sortSelect')
	const resetFilters = document.getElementById('resetFilters')
	const filtersToggle = document.getElementById('filtersToggle')
	const themeToggle = document.getElementById('themeToggle')
	const closeButton = document.querySelector('.modal-close')

	if (searchInput) searchInput.addEventListener('input', filterMembers)
	if (roleFilter) roleFilter.addEventListener('change', filterMembers)
	if (ageFilter) ageFilter.addEventListener('change', filterMembers)
	if (countryFilter) countryFilter.addEventListener('change', filterMembers)
	if (achievementFilter)
		achievementFilter.addEventListener('change', filterMembers)
	if (sortSelect) sortSelect.addEventListener('change', filterMembers)
	if (resetFilters) resetFilters.addEventListener('click', resetFilters)
	if (filtersToggle) filtersToggle.addEventListener('click', toggleFilters)
	if (themeToggle) themeToggle.addEventListener('click', toggleTheme)
	if (closeButton)
		closeButton.addEventListener('click', event => {
			event.stopPropagation()
			closeModal()
		})

	// Добавляем обработчик клавиши Escape для закрытия модального окна
	document.addEventListener('keydown', e => {
		if (e.key === 'Escape') {
			const modal = document.getElementById('userModal')
			if (!modal.classList.contains('hidden')) {
				closeModal()
			}
		}
	})

	// Добавляем обработчик клика вне модального окна
	document.addEventListener('click', event => {
		const modal = document.getElementById('userModal')
		if (modal && event.target === modal) {
			closeModal()
		}
	})
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
	// Загружаем тему
	loadTheme()

	// Добавляем обработчик для переключения темы
	const themeToggle = document.getElementById('themeToggle')
	if (themeToggle) {
		themeToggle.addEventListener('click', toggleTheme)
	}

	// Инициализируем остальные обработчики
	initializeEventListeners()
	fetchMembers()
})

// Функция для сохранения темы
function saveTheme(theme) {
	localStorage.setItem('theme', theme)
}

// Функция для загрузки темы
function loadTheme() {
	const savedTheme = localStorage.getItem('theme')
	if (savedTheme) {
		document.documentElement.classList.remove('light', 'dark')
		document.documentElement.classList.add(savedTheme)
		return savedTheme
	}

	// Если тема не сохранена, используем системные настройки
	const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
	const theme = prefersDark ? 'dark' : 'light'
	document.documentElement.classList.add(theme)
	saveTheme(theme)
	return theme
}

// Функция для переключения темы
function toggleTheme() {
	const currentTheme = document.documentElement.classList.contains('dark')
		? 'dark'
		: 'light'
	const newTheme = currentTheme === 'dark' ? 'light' : 'dark'

	document.documentElement.classList.remove(currentTheme)
	document.documentElement.classList.add(newTheme)
	saveTheme(newTheme)
}

document.getElementById('themeToggle').addEventListener('click', () => {
	const isDark = document.documentElement.classList.contains('dark')
	saveTheme(!isDark ? 'dark' : 'light')
})

loadTheme()

function toggleFilters() {
	const filtersPanel = document.getElementById('filtersPanel')
	const isHidden = filtersPanel.classList.contains('hidden')

	if (isHidden) {
		filtersPanel.classList.remove('hidden')
		setTimeout(() => {
			filtersPanel.classList.add('opacity-100', 'translate-y-0')
			filtersPanel.classList.remove('opacity-0', '-translate-y-2')
		}, 10)
	} else {
		filtersPanel.classList.add('opacity-0', '-translate-y-2')
		filtersPanel.classList.remove('opacity-100', 'translate-y-0')
		setTimeout(() => {
			filtersPanel.classList.add('hidden')
		}, 300)
	}
}
