function iceFound() {
  const ice = document.querySelector('.ice')

  ice.classList.add('found')

  setTimeout(() => {
    alert('Oei, das nou balen.')
  }, 500)
}

function bakkenAlert() {
  alert('Niet iedereen van Invictus speelt dit jaar mee met "Pirates of the Colosseum", maar dat weerhoudt ze niet van bakken vouwen. Deze mensen zijn aangegeven met een sterretje')
}

window.onload = function () {
  console.log('Deze site is natuurlijk alleen maar gemaakt voor mobiel. Als je daar wat van vind ben je uitgenodigd om met mij een bakje te komen vouwen in de vestingbar.')
}