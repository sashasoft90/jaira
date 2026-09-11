# Запись демонстрации в GIF. Требует vhs + ttyd + ffmpeg:
#   go install github.com/charmbracelet/vhs@latest
#   sudo apt install ttyd ffmpeg
#   vhs scripts/demo.tape        → docs/img/demo.gif
#
# Файл, а не запись рукой: когда поведение изменится, GIF
# переснимается одной командой, а не заново отыгрывается.

Output docs/img/demo.gif

Set Shell bash
Set FontSize 16
Set Width 1200
Set Height 800
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 25ms
Set PlaybackSpeed 1.0

Hide
Type "cd $(pwd) && go build -o /tmp/jaira ./cmd/jaira && clear" Enter
Sleep 5s
Show

Type "bash scripts/demo.sh"
Sleep 500ms
Enter
Sleep 60s
