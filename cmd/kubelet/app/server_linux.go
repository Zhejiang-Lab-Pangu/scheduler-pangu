/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"k8s.io/utils/inotify"

	libcontainercgroups "github.com/opencontainers/runc/libcontainer/cgroups"
)

func watchForLockfileContention(path string, done chan struct{}) error {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		klog.ErrorS(err, "Unable to create watcher for lockfile")
		return err
	}
	if err = watcher.AddWatch(path, inotify.InOpen|inotify.InDeleteSelf); err != nil {
		klog.ErrorS(err, "Unable to watch lockfile")
		watcher.Close()
		return err
	}
	go func() {
		select {
		case ev := <-watcher.Event:
			klog.InfoS("Inotify event", "event", ev)
		case err = <-watcher.Error:
			klog.ErrorS(err, "inotify watcher error")
		}
		close(done)
		watcher.Close()
	}()
	return nil
}

func isCgroup2UnifiedMode() bool {
	return libcontainercgroups.IsCgroup2UnifiedMode()
}
