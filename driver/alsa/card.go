package alsa

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	devPath     = "/dev"
	procPath    = "/proc"
	cardPath    = devPath + "/snd"
	modulesPath = procPath + "/asound/modules"
)

type Card struct {
	Path    string
	ID      string
	Name    string
	Driver  string
	Index   int
	Version Version

	fh   *os.File
	fd   uintptr
	info cardInfo
}

// Open a handle to a card by number.
func Open(index int) (*Card, error) {
	card := &Card{
		Path:  filepath.Join(cardPath, fmt.Sprintf("controlC%d", index)),
		Index: index,
	}

	var err error
	if card.fh, err = os.Open(card.Path); err != nil {
		return nil, err
	}
	card.fd = card.fh.Fd()

	if err = ioctl(card.fd, ioctlPointer(cmdRead, &card.Version, cmdControlVersion), &card.Version); err != nil {
		_ = card.fh.Close()
		return card, err
	}
	if err = ioctl(card.fd, ioctlPointer(cmdRead, &card.info, cmdControlCardInfo), &card.info); err != nil {
		_ = card.fh.Close()
		return card, err
	}

	card.ID = cstr(card.info.ID[:])
	card.Name = cstr(card.info.Name[:])
	card.Driver = cstr(card.info.Driver[:])
	return card, nil
}

// OpenDriver opens a handle to a card by driver name.
func OpenDriver(name string) (*Card, error) {
	f, err := os.Open(modulesPath)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.Split(strings.TrimSpace(s.Text()), " ")
		if len(line) < 2 {
			continue
		}
		if index, err := strconv.Atoi(line[0]); err == nil && strings.EqualFold(line[1], name) {
			_ = f.Close()
			return Open(index)
		}
	}

	return nil, fmt.Errorf("alsa: no card for driver %q was found", name)
}

// List all available cards.
func List() ([]*Card, error) {
	infos, err := os.ReadDir(cardPath)
	if err != nil {
		return nil, err
	}

	sort.Slice(infos, func(i, j int) bool {
		return strings.Compare(infos[i].Name(), infos[j].Name()) < 0
	})

	var cards []*Card
	for _, info := range infos {
		var index int
		if _, err = fmt.Sscanf(info.Name(), "controlC%d", &index); err != nil {
			continue
		}

		var card *Card
		if card, err = Open(index); err != nil {
			return cards, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (card *Card) Close() error {
	return card.fh.Close()
}

func (card *Card) Devices() ([]*Device, error) {
	var (
		devs []*Device
		next int32 = -1
		err  error
	)
	for {
		if err = ioctl(card.fd, ioctlPointer(cmdRead, &next, cmdControlPCMNextDevice), &next); err != nil {
			return nil, err
		} else if next == -1 {
			break // no more cards
		}

		var stream int32
		for ; stream < 2; stream++ {
			var (
				info = pcmInfo{
					Device: uint32(next),
					Card:   int32(card.Index),
					Stream: stream,
				}
				canPlay   = true
				canRecord = true
				ext       = "p"
			)
			if err = ioctl(card.fd, ioctlPointer(cmdRead|cmdWrite, &info, cmdControlPCMInfo), &info); err != nil {
				continue
			}

			if stream == 1 {
				canPlay = false
				canRecord = false
				ext = "c"
			}

			devs = append(devs, &Device{
				Type:      PCM,
				Path:      fmt.Sprintf("/dev/snd/pcmC%dD%d%s", card.Index, next, ext),
				CanPlay:   canPlay,
				CanRecord: canRecord,
				Index:     int(next),
				Name:      cstr(info.Name[:]),
				info:      info,
			})
		}
	}

	return devs, nil
}

type cardInfo struct {
	Card       int32
	_          int32
	ID         [16]byte
	Driver     [16]byte
	Name       [32]byte
	LongName   [80]byte
	_          [16]byte
	MixerName  [80]byte
	Components [128]byte
}

type pcmInfo struct {
	Device          uint32   /* RO/WR (control): device number */
	Subdevice       uint32   /* RO/WR (control): subdevice number */
	Stream          int32    /* RO/WR (control): stream direction */
	Card            int32    /* R: card number */
	ID              [64]byte /* ID (user selectable) */
	Name            [80]byte /* name of this device */
	Subname         [32]byte /* subdevice name */
	DevClass        int32    /* SNDRV_PCM_CLASS_* */
	DevSubclass     int32    /* SNDRV_PCM_SUBCLASS_* */
	SubdevicesCount uint32
	SubdevicesAvail uint32
	SyncID          [16]byte /* hardware synchronization ID */
	_               [64]byte /* reserved for future... */
}
