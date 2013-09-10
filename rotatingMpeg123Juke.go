package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "os/exec"
)

/**
 * Un petit juke-box en go :
 *  - Ce programme lit les mp3s listés dans le fichier playlist.txt
 *  - La lecture est séquentielle
 *  - Suite à l'arrêt du programme, la lecture reprend au fichier suivant
 *    de la playlist
 */

// Charge le fichier spécifié 
// et retourne un slice de ses lignes
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

// Ecrit la ligne dans le fichier spécifié 
func writeLine(line string, path string) error {
  file, err := os.Create(path)
  if err != nil {
    return err
  }
  defer file.Close()

  w := bufio.NewWriter(file)
  fmt.Fprintln(w, line)

  return w.Flush()
}

func launchMP3(mp3File string) {
    app := "mpg123"

    cmd := exec.Command(app, mp3File)
    out, err := cmd.Output()

    if err != nil {
        println(err.Error())
        return
    }

    print(string(out))
}

func updateNextTrack(nextTrack string) {
      errw := writeLine(nextTrack, "nexttrack.txt")
      if errw != nil {
        log.Fatalf("writeLine: %s", errw)
      }
}

func main() {
    // Recherche du playlist.txt
    playList, err := readLines("playlist.txt")
    if err != nil {
      log.Fatalf("readLines: %s", err)
    }

    // Recherche du nexttrack.txt
    nextTracks, errT := readLines("nexttrack.txt")
    if errT != nil {
      log.Printf("readLines: %s", errT)
      // On essaye de créer le fichier s'il n'existe pas
      nextTracks = make([]string,1,1)
      nextTracks[0] = playList[0]
      updateNextTrack(nextTracks[0])
    }
    nextTrack := nextTracks[0]
    nbLines := len(playList)

    // Recherche de la piste suivante
    startPos := -1

    for i, line := range playList {
      if line == nextTrack {
        startPos = i
      }
    }

    if startPos == -1 {
      log.Printf("Impossible de trouver l'index de %s", nextTrack)

      // On repart à 0
      nextTrack = playList[0]
      nextTracks[0] = nextTrack
      updateNextTrack(nextTrack)
      startPos = 0
    }

    // Boucle sur la playlist à partir de nextTrack
    for {
      nextMP3 := playList[startPos]
      startPos++

      // Mise à jour de nextTrack
      if startPos == nbLines {
        startPos = 0
      }
      nextTrack = playList[startPos]
      nextTracks[0] = nextTrack
      updateNextTrack(nextTrack)

      // Lancement du mp3 qui va bien
      fmt.Printf("Lancement de %s\n", nextMP3)
      launchMP3(nextMP3)
    }
}
