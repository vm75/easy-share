package server

import (
	"fmt"
	"os"
	"strings"
)

const globalStaticConfig string = `
   server role = standalone server
   log file = /dev/stdout
   max log size = 50
   dns proxy = no

   # Other settings
   pam password change = yes
   map to guest = bad user
   usershare allow guests = yes
   create mask = 0664
   force create mode = 0664
   directory mask = 0775
   force directory mode = 0775
   force user = smbuser
   force group = smb
   follow symlinks = yes
   load printers = no
   printing = bsd
   printcap name = /dev/null
   disable spoolss = yes
   strict locking = no
   aio read size = 0
   aio write size = 0
   vfs objects = catia fruit recycle streams_xattr
   recycle:keeptree = yes
   recycle:maxsize = 0
   recycle:repository = .deleted
   recycle:touch_mtime = yes
   recycle:versions = yes
   recycle:noversions = *.tmp,*.temp,*.o,*.obj,*.TMP,*.TEMP
   recycle:exclude = *.tmp,*.temp,*.o,*.obj,*.TMP,*.TEMP,*.bak,thumb.db,Thumbs.db,*.bak,*~,*.swp,*Zone.Identifier
	 recycle:exclude_dir = /trash,/recycle,/tmp,/temp,/TMP,/TEMP,.recycle

   # Security
   client ipc max protocol = SMB3
   client ipc min protocol = SMB2_10
   client max protocol = SMB3
   client min protocol = SMB2_10
   server max protocol = SMB3
   server min protocol = SMB2_10

   # Time Machine
   fruit:delete_empty_adfiles = yes
   fruit:time machine = yes
   fruit:veto_appledouble = no
   fruit:wipe_intentionally_left_blank_rfork = yes
`

const recycleSettings string = `
   recycle:keeptree = yes
   recycle:maxsize = 0
   recycle:repository = .deleted
   recycle:touch_mtime = yes
   recycle:versions = yes
   recycle:noversions = *.tmp,*.temp,*.o,*.obj,*.TMP,*.TEMP
   recycle:exclude = *.tmp,*.temp,*.o,*.obj,*.TMP,*.TEMP,*.bak,thumb.db,Thumbs.db,*.bak,*~,*.swp,*Zone.Identifier
	 recycle:exclude_dir = /trash,/recycle,/tmp,/temp,/TMP,/TEMP,.recycle
`

const vetoFiles string = `
   veto files = /.apdisk/.DS_Store/.TemporaryItems/.Trashes/desktop.ini/ehthumbs.db/Network Trash Folder/Temporary Items/Thumbs.db/*Zone.Identifier/
   delete veto files = yes
`

func UpdateSambaConfig(globalConfig SambaGlobalConfig, shareConfigs map[string]SambaShareConfig) error {
	var sb strings.Builder
	sb.WriteString("[global]\n")
	sb.WriteString(fmt.Sprintf("workgroup = %s\n", globalConfig.Workgroup))
	sb.WriteString(fmt.Sprintf("server string = %s\n", globalConfig.ServerString))
	sb.WriteString(fmt.Sprintf("guest user = %s\n", globalConfig.GuestUser))
	if len(globalConfig.AllowedHosts) > 0 {
		sb.WriteString(fmt.Sprintf("allowed hosts = %s\n", strings.Join(globalConfig.AllowedHosts, " ")))
	}
	if globalConfig.EnableRecycle {
		sb.WriteString(fmt.Sprintf("vfs objects = %s\n", strings.Join(append(globalConfig.VfsObjects, "recycle"), " ")))
	}
	if globalConfig.EnableRecycle {
		sb.WriteString(recycleSettings)
	}
	sb.WriteString(globalStaticConfig)
	sb.WriteString("\n")

	for _, shareConfig := range shareConfigs {
		sb.WriteString(fmt.Sprintf("[%s]\n", shareConfig.Name))
		if shareConfig.Comment != "" {
			sb.WriteString(fmt.Sprintf("comment = %s\n", shareConfig.Comment))
		}
		sb.WriteString(fmt.Sprintf("path = %s\n", shareConfig.Path))
		sb.WriteString(fmt.Sprintf("browsable = %t\n", shareConfig.Browsable))
		sb.WriteString(fmt.Sprintf("writable = %t\n", shareConfig.Writable))
		sb.WriteString(fmt.Sprintf("guest ok = %t\n", shareConfig.GuestOk))
		if len(shareConfig.Users) > 0 {
			sb.WriteString(fmt.Sprintf("users = %s\n", strings.Join(shareConfig.Users, " ")))
		}
		if len(shareConfig.Admins) > 0 {
			sb.WriteString(fmt.Sprintf("admins = %s\n", strings.Join(shareConfig.Admins, " ")))
		}
		if len(shareConfig.Writelist) > 0 {
			sb.WriteString(fmt.Sprintf("writelist = %s\n", strings.Join(shareConfig.Writelist, " ")))
		}
		if shareConfig.Veto {
			sb.WriteString(vetoFiles)
		}
		if shareConfig.CatiaMappings != "" {
			sb.WriteString(fmt.Sprintf("catia:mappings = %s\n", shareConfig.CatiaMappings))
		}
		if shareConfig.CreateMask != "" {
			sb.WriteString(fmt.Sprintf("create mask = %s\n", shareConfig.CreateMask))
		}
		if shareConfig.CustomOptions != "" {
			sb.WriteString(fmt.Sprintf("custom options = %s\n", shareConfig.CustomOptions))
		}
		sb.WriteString("\n")
	}

	err := os.WriteFile("/etc/samba/smb.conf", []byte(sb.String()), 0644)

	return err
}

func UpdateNfsConfig(shareConfigs map[string][]NfsShareConfig) error {
	var sb strings.Builder

	for path, hosts := range shareConfigs {
		// skip if path does not exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		sb.WriteString(path)
		for _, host := range hosts {
			var writable string
			if host.Writable {
				writable = "rw"
			} else {
				writable = "ro"
			}
			sb.WriteString(fmt.Sprintf(" %s(%s", host.Host, writable))
			if !host.Secure {
				sb.WriteString(",insecure")
			}
			if !host.Sync {
				sb.WriteString(",async")
			}
			if host.Mapping != "" {
				sb.WriteString("," + host.Mapping)
			}
			if host.Anonuid != 0 {
				sb.WriteString(fmt.Sprintf(",anonuid=%d", host.Anonuid))
			}
			if host.Anongid != 0 {
				sb.WriteString(fmt.Sprintf(",anongid=%d", host.Anongid))
			}
			if host.CustomOptions != "" {
				sb.WriteString("," + host.CustomOptions)
			}
			sb.WriteString(")")
		}
		sb.WriteString("\n")
	}

	err := os.WriteFile("/etc/exports", []byte(sb.String()), 0644)

	return err
}
