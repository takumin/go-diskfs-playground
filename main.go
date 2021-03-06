package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/partition/gpt"
)

func saveFile(fs filesystem.FileSystem, dst, src string) error {
	file, err := fs.OpenFile(dst, os.O_CREATE|os.O_RDWR)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func getFileSize(f string) int64 {
	file, err := os.OpenFile(f, os.O_RDONLY, 0600)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}
	if stat.IsDir() {
		log.Panic("require kernel image file path")
	}

	return stat.Size()
}

func main() {
	var (
		diskImg       string
		kernImg       string
		initImg       string
		rootImg       string
		cmdline       string
		metaData      string
		userData      string
		networkConfig string
	)

	flag.StringVar(&diskImg, "disk", "/tmp/disk.img", "disk image file path")
	flag.StringVar(&kernImg, "kernel", "", "kernel image file path")
	flag.StringVar(&initImg, "initrd", "", "initrd image file path")
	flag.StringVar(&rootImg, "rootfs", "", "rootfs image file path")
	flag.StringVar(&cmdline, "cmdline", "", "kernel boot arguments")
	flag.StringVar(&metaData, "metaData", "", "cloud-init nocloud meta-data file path")
	flag.StringVar(&userData, "userData", "", "cloud-init nocloud user-data file path")
	flag.StringVar(&networkConfig, "networkConfig", "", "cloud-init nocloud network-config file path")

	flag.Parse()

	if diskImg == "" || kernImg == "" {
		flag.PrintDefaults()
		return
	}

	imgSize := getFileSize(kernImg)
	if initImg != "" {
		imgSize += getFileSize(initImg)
	}
	if rootImg != "" {
		imgSize += getFileSize(rootImg)
	}
	if metaData != "" {
		imgSize += getFileSize(metaData)
	}
	if userData != "" {
		imgSize += getFileSize(userData)
	}
	if networkConfig != "" {
		imgSize += getFileSize(networkConfig)
	}
	if imgSize%(1024*1024) != 0 {
		imgSize += (1024 * 1024) - (imgSize % (1024 * 1024))
	}

	var (
		espSize          int64 = imgSize + 4*1024*1024
		diskSize         int64 = espSize + 4*1024*1024
		blkSize          int64 = 512
		partitionStart   int64 = 2048
		partitionSectors int64 = espSize / blkSize
		partitionEnd     int64 = partitionSectors - partitionStart + 1
	)

	img, err := diskfs.Create(diskImg, diskSize, diskfs.Raw)
	if err != nil {
		log.Panic(err)
	}

	table := &gpt.Table{
		LogicalSectorSize:  512,
		PhysicalSectorSize: 512,
		ProtectiveMBR:      true,
		Partitions: []*gpt.Partition{
			&gpt.Partition{
				Name:  "ESP",
				Type:  gpt.EFISystemPartition,
				Start: uint64(partitionStart),
				End:   uint64(partitionEnd),
			},
		},
	}

	err = img.Partition(table)
	if err != nil {
		log.Panic(err)
	}

	spec := disk.FilesystemSpec{
		Partition: 1,
		FSType:    filesystem.TypeFat32,
	}

	fs, err := img.CreateFilesystem(spec)
	if err != nil {
		log.Panic(err)
	}

	err = fs.Mkdir("/EFI/BOOT")
	if err != nil {
		log.Panic(err)
	}

	if err = saveFile(fs, "/EFI/BOOT/BOOTX64.EFI", kernImg); err != nil {
		log.Panic(err)
	}

	log.Println("Finish!")
}
