	"path"
	"github.com/m3db/m3/src/x/lockfile"
	filePathPrefixLockFile           = ".lock"
	newFileMode, err := cfg.Filesystem.ParseNewFileMode()
	if err != nil {
		logger.Fatalf("could not parse new file mode: %v", err)
	}

	newDirectoryMode, err := cfg.Filesystem.ParseNewDirectoryMode()
	if err != nil {
		logger.Fatalf("could not parse new directory mode: %v", err)
	}

	// Obtain a lock on `filePathPrefix`, or exit if another process already has it.
	// The lock consists of a lock file (on the file system) and a lock in memory.
	// When the process exits gracefully, both the lock file and the lock will be removed.
	// If the process exits ungracefully, only the lock in memory will be removed, the lock
	// file will remain on the file system. When a dbnode starts after an ungracefully stop,
	// it will be able to acquire the lock despite the fact the the lock file exists.
	lockPath := path.Join(cfg.Filesystem.FilePathPrefixOrDefault(), filePathPrefixLockFile)
	fslock, err := lockfile.CreateAndAcquire(lockPath, newDirectoryMode)
	if err != nil {
		logger.Fatalf("could not acquire lock on %s: %v", lockPath, err)
	}
	defer fslock.Release()
