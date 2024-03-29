package window

import (
	"log"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/zapomnij/firecraft/pkg/downloader"
)

var (
	lpf *LauncherProfiles
	vm  *downloader.VersionManifest
)

type FWindow struct {
	Window *widgets.QWidget

	container *widgets.QWidget
	layout    *widgets.QGridLayout

	ProgressBar       *widgets.QProgressBar
	ProgressBarStatus *widgets.QLabel
	gameLogger        *widgets.QPlainTextEdit
	notebook          *widgets.QTabWidget

	ms *MSTab

	bottomBar       *widgets.QWidget
	bottomBarLayout *widgets.QHBoxLayout
	playBt          *widgets.QPushButton

	userBox    *widgets.QWidget
	userLay    *widgets.QGridLayout
	usernameTv *widgets.QLineEdit

	profilesBox      *widgets.QWidget
	profilesLay      *widgets.QGridLayout
	profilesSelector *widgets.QComboBox
	editProfile      *widgets.QPushButton
}

func NewFWindow() *FWindow {
	var err error
	lpf, err = loadProfiles()
	if err != nil {
		log.Fatalln(err)
	}

	vm, err = downloader.GetVersionManifest()
	if err != nil {
		log.Fatalln(err)
	}

	var this = FWindow{}
	this.Window = widgets.NewQWidget(nil, 0)
	this.Window.SetWindowTitle(core.QCoreApplication_ApplicationName())

	this.container = widgets.NewQWidget(this.Window, 0)
	this.layout = widgets.NewQGridLayout(this.container)
	this.layout.SetSpacing(0)
	this.layout.SetContentsMargins(0, 0, 0, 0)

	this.notebook = widgets.NewQTabWidget(this.container)
	this.gameLogger = widgets.NewQPlainTextEdit(this.notebook)
	this.gameLogger.SetReadOnly(true)
	this.notebook.AddTab(this.gameLogger, "Game logs")
	this.layout.AddWidget(this.notebook)

	this.ProgressBar = widgets.NewQProgressBar(this.container)
	progressBarLayout := widgets.NewQVBoxLayout2(this.ProgressBar)
	progressBarLayout.SetContentsMargins(0, 0, 0, 0)
	this.ProgressBarStatus = widgets.NewQLabel2("status", this.ProgressBar, 0)
	progressBarLayout.AddWidget(this.ProgressBarStatus, 0, core.Qt__AlignCenter)
	this.ProgressBar.SetLayout(progressBarLayout)
	this.ProgressBar.SetVisible(false)
	this.ProgressBar.SetTextVisible(false)
	this.layout.AddWidget(this.ProgressBar)

	this.bottomBar = widgets.NewQWidget(this.container, 0)
	this.bottomBar.SetFixedHeight(70)
	this.bottomBarLayout = widgets.NewQHBoxLayout()
	this.bottomBarLayout.SetContentsMargins(0, 0, 0, 0)
	this.bottomBar.SetLayout(this.bottomBarLayout)

	this.profilesBox = widgets.NewQWidget(this.bottomBar, 0)
	this.profilesLay = widgets.NewQGridLayout(this.profilesBox)
	this.profilesBox.SetLayout(this.profilesLay)
	this.editProfile = widgets.NewQPushButton2("Edit profile", this.profilesBox)
	this.editProfile.ConnectClicked(this.editProfileHandle)
	this.profilesSelector = widgets.NewQComboBox(this.profilesBox)
	this.reloadProfileSelector(lpf.PreviousProfile)
	this.profilesSelector.ConnectCurrentTextChanged(this.updatePreviousProfile)

	this.profilesLay.AddWidget(this.profilesSelector)
	this.profilesLay.AddWidget(this.editProfile)

	this.playBt = widgets.NewQPushButton2("Play", this.bottomBar)
	this.playBt.ConnectClicked(func(checked bool) {
		this.playBt.SetEnabled(false)
		go this.Launch()
	})
	this.playBt.SetFixedHeight(60)
	this.playBt.SetFixedWidth(300)

	this.userBox = widgets.NewQWidget(this.bottomBar, 0)
	this.userLay = widgets.NewQGridLayout(this.userBox)
	this.userLay.AddWidget(widgets.NewQLabel2("Username", this.userBox, 0))
	this.usernameTv = widgets.NewQLineEdit2(lpf.AuthenticationDatabase.Username, this.userBox)
	this.usernameTv.ConnectTextChanged(this.saveUsername)
	this.userLay.AddWidget(this.usernameTv)
	this.userBox.SetLayout(this.userLay)

	this.bottomBarLayout.AddWidget(this.profilesBox, 0, core.Qt__AlignLeft)
	this.bottomBarLayout.AddWidget(this.playBt, 0, core.Qt__AlignHCenter)
	this.bottomBarLayout.AddWidget(this.userBox, 0, core.Qt__AlignRight)

	this.layout.AddWidget(this.bottomBar)

	this.ms = NewMSTab(this.notebook, &this)
	this.notebook.AddTab(this.ms.widget, "MS Authentication")

	this.container.SetLayout(this.layout)
	this.Window.SetLayout(this.layout)

	if downloader.OperatingSystem == "osx" {
		this.macOSFix()
	}

	return &this
}

func (fw *FWindow) updatePreviousProfile(text string) {
	if text != "New profile" {
		lpf.PreviousProfile = text
		_ = lpf.Save()
	}
}

func (fw *FWindow) editProfileHandle(_ bool) {
	epw := NewEditProfileWindow(fw)
	epw.Window.Resize2(300, 200)
	epw.Window.Show()
}

func (fw *FWindow) saveUsername(text string) {
	lpf.AuthenticationDatabase.Username = text
	_ = lpf.Save()
}

func (fw *FWindow) reloadProfileSelector(set string) {
	fw.profilesSelector.Clear()

	i := 0
	for k := range lpf.Profiles {
		fw.profilesSelector.AddItem(k, core.NewQVariant())
		if k == set {
			fw.profilesSelector.SetCurrentIndex(i)
		}
		i++
	}

	fw.profilesSelector.AddItem("New profile", core.NewQVariant())
}

func (fw *FWindow) updateProgressBar(inc int, status string) {
	fw.ProgressBar.SetValue(fw.ProgressBar.Value() + inc)
	fw.ProgressBarStatus.SetText(status)
}

func (fw *FWindow) macOSFix() {
	fw.gameLogger.SetPlainText("Information for macOS users:\n\nYou should use JVM from homebrew")
	fw.userLay.SetContentsMargins(5, 5, 5, 5)
	fw.profilesLay.SetContentsMargins(5, 3, 5, 5)
}
