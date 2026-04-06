import { ref, onMounted } from 'vue'

const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID

export function useGoogleAuth() {
  const isGoogleLoaded = ref(false)

  // Inicializar o Google Identity Services
  const initGoogleAuth = () => {
    return new Promise((resolve) => {
      // Verificar se o script já foi carregado
      if (window.google) {
        isGoogleLoaded.value = true
        resolve(true)
        return
      }

      // Esperar o script carregar
      const checkGoogle = setInterval(() => {
        if (window.google) {
          clearInterval(checkGoogle)
          isGoogleLoaded.value = true
          resolve(true)
        }
      }, 100)

      // Timeout de 10 segundos
      setTimeout(() => {
        clearInterval(checkGoogle)
        resolve(false)
      }, 10000)
    })
  }

  // Renderizar o botão do Google
  const renderGoogleButton = (elementId, callback) => {
    if (!window.google || !GOOGLE_CLIENT_ID) {
      console.error('Google Identity Services não carregado ou Client ID não configurado')
      return
    }

    window.google.accounts.id.initialize({
      client_id: GOOGLE_CLIENT_ID,
      callback: callback,
      auto_select: false,
      cancel_on_tap_outside: true,
    })

    window.google.accounts.id.renderButton(
      document.getElementById(elementId),
      {
        theme: 'filled_black',
        size: 'large',
        width: '100%',
        text: 'signin_with',
        shape: 'rectangular',
        logo_alignment: 'left',
      }
    )
  }

  // Login com popup
  const loginWithPopup = () => {
    return new Promise((resolve, reject) => {
      if (!window.google) {
        reject(new Error('Google Identity Services não carregado'))
        return
      }

      window.google.accounts.id.initialize({
        client_id: GOOGLE_CLIENT_ID,
        callback: (response) => {
          resolve(response.credential)
        },
      })

      window.google.accounts.id.prompt((notification) => {
        if (notification.isNotDisplayed() || notification.isSkippedMoment()) {
          // One Tap não foi exibido, tentar método alternativo
          reject(new Error('One Tap não disponível'))
        }
      })
    })
  }

  return {
    isGoogleLoaded,
    initGoogleAuth,
    renderGoogleButton,
    loginWithPopup,
  }
}
