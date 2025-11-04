; TF-Engine 2.0 NSIS Installer Script
; Build command: makensis tf-engine-installer.nsi

!include "MUI2.nsh"
!include "LogicLib.nsh"

; Application metadata
!define APP_NAME "TF-Engine"
!define APP_VERSION "1.0.0"
!define APP_PUBLISHER "TF Systems"
!define APP_EXE "tf-engine.exe"
!define UNINSTALL_EXE "Uninstall.exe"

Name "${APP_NAME} ${APP_VERSION}"
OutFile "TFEngine-Setup-${APP_VERSION}.exe"
InstallDir "$PROGRAMFILES64\${APP_NAME}"
InstallDirRegKey HKLM "Software\${APP_NAME}" "InstallPath"

RequestExecutionLevel admin

; Modern UI Configuration
; Icon configuration (commented out - uncomment when icon.ico is available)
; !define MUI_ICON "assets\icon.ico"
; !define MUI_UNICON "assets\icon.ico"
!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_RUN "$INSTDIR\${APP_EXE}"
!define MUI_FINISHPAGE_RUN_TEXT "Launch TF-Engine"

; Installer pages
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "LICENSE.txt"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

; Version check function
Function .onInit
  ReadRegStr $0 HKLM "Software\${APP_NAME}" "Version"
  ${If} $0 != ""
    MessageBox MB_YESNO|MB_ICONQUESTION "TF-Engine is already installed (version $0). Continue with installation?" IDYES +2
    Abort
  ${EndIf}
FunctionEnd

; Main installation section
Section "Install"
  SetOutPath "$INSTDIR"

  ; Core files
  File "dist\${APP_EXE}"
  File "dist\policy.v1.json"
  File "LICENSE.txt"
  File "dist\README.md"

  ; Create data directory for trade history
  CreateDirectory "$INSTDIR\data"
  CreateDirectory "$INSTDIR\data\ui"

  ; Create feature flags file with Phase 2 features disabled
  FileOpen $0 "$INSTDIR\feature.flags.json" w
  FileWrite $0 '{"trade_management":false,"sample_data_generator":false,"vimium_mode":false,"advanced_analytics":false}'
  FileClose $0

  ; Registry keys
  WriteRegStr HKLM "Software\${APP_NAME}" "InstallPath" "$INSTDIR"
  WriteRegStr HKLM "Software\${APP_NAME}" "Version" "${APP_VERSION}"

  ; Uninstaller
  WriteUninstaller "$INSTDIR\${UNINSTALL_EXE}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayName" "${APP_NAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "UninstallString" "$INSTDIR\${UNINSTALL_EXE}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "Publisher" "${APP_PUBLISHER}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayVersion" "${APP_VERSION}"
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "EstimatedSize" 4096

  ; Start Menu shortcuts
  CreateDirectory "$SMPROGRAMS\${APP_NAME}"
  CreateShortcut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" "$INSTDIR\${APP_EXE}" "" "$INSTDIR\${APP_EXE}" 0
  CreateShortcut "$SMPROGRAMS\${APP_NAME}\Uninstall.lnk" "$INSTDIR\${UNINSTALL_EXE}" "" "$INSTDIR\${UNINSTALL_EXE}" 0

  ; Desktop shortcut (optional)
  MessageBox MB_YESNO "Create desktop shortcut?" IDNO +2
  CreateShortcut "$DESKTOP\${APP_NAME}.lnk" "$INSTDIR\${APP_EXE}" "" "$INSTDIR\${APP_EXE}" 0
SectionEnd

; Uninstaller
Section "Uninstall"
  ; Remove files
  Delete "$INSTDIR\${APP_EXE}"
  Delete "$INSTDIR\policy.v1.json"
  Delete "$INSTDIR\LICENSE.txt"
  Delete "$INSTDIR\README.md"
  Delete "$INSTDIR\feature.flags.json"
  Delete "$INSTDIR\${UNINSTALL_EXE}"

  ; Prompt to preserve user data
  MessageBox MB_YESNO "Delete trade history and settings?$\n(This cannot be undone)" IDYES delete_data
  Goto skip_data
  delete_data:
    RMDir /r "$INSTDIR\data"
  skip_data:

  ; Remove shortcuts
  Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
  Delete "$SMPROGRAMS\${APP_NAME}\Uninstall.lnk"
  RMDir "$SMPROGRAMS\${APP_NAME}"
  Delete "$DESKTOP\${APP_NAME}.lnk"

  ; Remove registry keys
  DeleteRegKey HKLM "Software\${APP_NAME}"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"

  ; Remove install directory (if empty)
  RMDir "$INSTDIR"
SectionEnd
