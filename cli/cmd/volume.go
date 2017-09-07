package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	log "github.com/Sirupsen/logrus"
	"github.com/gluster/glusterd2/pkg/api"
)

const (
	helpVolumeCmd       = "Gluster Volume Management"
	helpVolumeCreateCmd = "Create a Gluster Volume"
	helpVolumeStartCmd  = "Start a Gluster Volume"
	helpVolumeStopCmd   = "Stop a Gluster Volume"
	helpVolumeDeleteCmd = "Delete a Gluster Volume"
	helpVolumeGetCmd    = "Get Gluster Volume Options"
	helpVolumeSetCmd    = "Set a Gluster Volume Option"
	helpVolumeResetCmd  = "Reset a Gluster Volume Option"
	helpVolumeInfoCmd   = "Get Gluster Volume Info"
	helpVolumeListCmd   = "List all Gluster Volumes"
	helpVolumeStatusCmd = "Get Gluster Volume Status"
)

var (
	// Create Command Flags
	flagCreateCmdStripeCount       int
	flagCreateCmdReplicaCount      int
	flagCreateCmdDisperseCount     int
	flagCreateCmdDisperseDataCount int
	flagCreateCmdRedundancyCount   int
	flagCreateCmdTransport         string
	flagCreateCmdForce             bool

	// Start Command Flags
	flagStartCmdForce bool

	// Stop Command Flags
	flagStopCmdForce bool
)

func init() {
	// Volume Create
	volumeCreateCmd.Flags().IntVarP(&flagCreateCmdStripeCount, "stripe", "", 0, "Stripe Count")
	volumeCreateCmd.Flags().IntVarP(&flagCreateCmdReplicaCount, "replica", "", 0, "Replica Count")
	volumeCreateCmd.Flags().IntVarP(&flagCreateCmdDisperseCount, "disperse", "", 0, "Disperse Count")
	volumeCreateCmd.Flags().IntVarP(&flagCreateCmdDisperseDataCount, "disperse-data", "", 0, "Disperse Data Count")
	volumeCreateCmd.Flags().IntVarP(&flagCreateCmdRedundancyCount, "redundancy", "", 0, "Redundancy Count")
	volumeCreateCmd.Flags().StringVarP(&flagCreateCmdTransport, "transport", "", "tcp", "Transport")
	volumeCreateCmd.Flags().BoolVarP(&flagCreateCmdForce, "force", "f", false, "Force")
	volumeCmd.AddCommand(volumeCreateCmd)

	// Volume Start
	volumeStartCmd.Flags().BoolVarP(&flagStartCmdForce, "force", "f", false, "Force")
	volumeCmd.AddCommand(volumeStartCmd)

	// Volume Stop
	volumeStopCmd.Flags().BoolVarP(&flagStopCmdForce, "force", "f", false, "Force")
	volumeCmd.AddCommand(volumeStopCmd)

	// Volume Delete
	volumeCmd.AddCommand(volumeDeleteCmd)

	volumeCmd.AddCommand(volumeGetCmd)
	volumeCmd.AddCommand(volumeSetCmd)
	volumeCmd.AddCommand(volumeResetCmd)
	volumeCmd.AddCommand(volumeInfoCmd)
	volumeCmd.AddCommand(volumeStatusCmd)
	RootCmd.AddCommand(volumeCmd)
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: helpVolumeCmd,
}

var volumeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: helpVolumeCreateCmd,
	Run: func(cmd *cobra.Command, args []string) {
		validateNArgs(cmd, 2, 0)
		volname := cmd.Flags().Args()[0]
		bricks := cmd.Flags().Args()[1:]
		fmt.Println("CREATE:", volname, bricks)
	},
}

var volumeStartCmd = &cobra.Command{
	Use:   "start [flags] <VOLNAME>",
	Short: helpVolumeStartCmd,
	Run: func(cmd *cobra.Command, args []string) {
		validateNArgs(cmd, 1, 1)
		volname := cmd.Flags().Args()[0]
		err := client.VolumeStart(volname)
		if err != nil {
			log.WithField("volume", volname).Println("volume start failed")
			failure(fmt.Sprintf("volume start failed with: %s", err.Error()), 1)
		}
		fmt.Printf("Volume %s started successfully\n", volname)
	},
}

var volumeStopCmd = &cobra.Command{
	Use:   "stop [flags] <VOLNAME>",
	Short: helpVolumeStopCmd,
	Run: func(cmd *cobra.Command, args []string) {
		validateNArgs(cmd, 1, 1)
		volname := cmd.Flags().Args()[0]
		err := client.VolumeStop(volname)
		if err != nil {
			log.WithField("volume", volname).Println("volume stop failed")
			failure(fmt.Sprintf("volume stop failed with: %s", err.Error()), 1)
		}
		fmt.Printf("Volume %s stopped successfully\n", volname)
	},
}

var volumeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: helpVolumeDeleteCmd,
	Run: func(cmd *cobra.Command, args []string) {
		validateNArgs(cmd, 1, 1)
		volname := cmd.Flags().Args()[0]
		fmt.Println("DELETE:", volname)
	},
}

var volumeGetCmd = &cobra.Command{
	Use:   "get",
	Short: helpVolumeGetCmd,
	Run: func(cmd *cobra.Command, args []string) {
		validateNArgs(cmd, 1, 2)
		volname := cmd.Flags().Args()[0]
		fmt.Println("GET:", volname)
	},
}

var volumeSetCmd = &cobra.Command{
	Use:   "set",
	Short: helpVolumeSetCmd,
	Run: func(cmd *cobra.Command, args []string) {
		validateNArgs(cmd, 3, 3)
		volname := cmd.Flags().Args()[0]
		fmt.Println("SET:", volname)
	},
}

var volumeResetCmd = &cobra.Command{
	Use:   "reset",
	Short: helpVolumeResetCmd,
	Run: func(cmd *cobra.Command, args []string) {
		validateNArgs(cmd, 1, 2)
		volname := cmd.Flags().Args()[0]
		fmt.Println("RESET:", volname)
	},
}

func volumeInfoHandler2(cmd *cobra.Command)(api.Volinfo, error) {
	var vols api.Volinfo
	validateNArgs(cmd, 0, 1)
	volname := ""
	if len(cmd.Flags().Args()) > 0 {
		volname = cmd.Flags().Args()[0]
	}
	if volname == ""{
		vols_, err := client.VolumeAll()
		if err != nil {
			log.Println("No volumes present")
			failure(fmt.Sprintf("Error getting Volumes list %s", err.Error()), 1)
		}
		fmt.Println(vols_)
	}else{
		vols, _ = client.Volumes(volname)
	}
	return vols, nil
}

var volumeInfoCmd = &cobra.Command{
	Use:   "info",
	Short: helpVolumeInfoCmd,
	Run: func(cmd *cobra.Command, args []string) {
		vols, err := volumeInfoHandler2(cmd)
		fmt.Println(vols)
		fmt.Println(err)
		// for i, vol := range vols {
		// 		fmt.Println("Volume Name: ", vol.Name)
		// }
	},
}

var volumeListCmd = &cobra.Command{
	Use:   "list",
	Short: helpVolumeListCmd,
	Run: func(cmd *cobra.Command, args []string) {
		volumeInfoHandler2(cmd)
	},
}

var volumeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: helpVolumeStatusCmd,
	Run: func(cmd *cobra.Command, args []string) {
		validateNArgs(cmd, 0, 1)
		volname := "all"
		if len(cmd.Flags().Args()) > 0 {
			volname = cmd.Flags().Args()[0]
		}
		fmt.Println("STATUS:", volname)
	},
}
