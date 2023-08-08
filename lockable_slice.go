package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

type _LockableSlice struct {
	slice  []interface{}
	locked bool
	index  int
}

func _NewLockableSlice() *_LockableSlice {
	return &_LockableSlice{
		slice: []interface{}{},
	}
}

func (this *_LockableSlice) _RequireNotLocked() {
	if this.locked {
		panic(errLockedSlice)
	}
}

func (this *_LockableSlice) _SetLocked(locked bool) *_LockableSlice { // nolint
	this.locked = locked
	return this
}

func (this *_LockableSlice) _SetSlice(slice []interface{}) *_LockableSlice { //nolint
	this._RequireNotLocked()
	this.slice = slice
	this.index = 0
	return this
}

func (this *_LockableSlice) _Push(items ...interface{}) *_LockableSlice {
	this._RequireNotLocked()
	this.slice = append(this.slice, items...)
	return this
}

func (this *_LockableSlice) _Clear() *_LockableSlice { //nolint
	this._RequireNotLocked()
	this.slice = []interface{}{}
	return this
}

func (this *_LockableSlice) _Get(index int) interface{} { //nolint
	return this.slice[index]
}

func (this *_LockableSlice) _Set(index int, item interface{}) *_LockableSlice { //nolint
	this._RequireNotLocked()

	if len(this.slice) == index {
		this.slice = append(this.slice, item)
	} else {
		this.slice[index] = item
	}

	return this
}

func (this *_LockableSlice) _SetIfAbsent(index int, item interface{}) *_LockableSlice { //nolint
	this._RequireNotLocked()
	if len(this.slice) == index || this.slice[index] == nil {
		this._Set(index, item)
	}
	return this
}

func (this *_LockableSlice) _GetNext() interface{} { //nolint
	return this._Get(this._Advance())
}

func (this *_LockableSlice) _GetCurrent() interface{} { //nolint
	return this._Get(this.index)
}

func (this *_LockableSlice) _Advance() int { //nolint
	index := this.index
	if len(this.slice) != 0 {
		this.index = (this.index + 1) % len(this.slice)
	}
	return index
}

func (this *_LockableSlice) _IsEmpty() bool { //nolint
	return len(this.slice) == 0
}

func (this *_LockableSlice) _Length() int { //nolint
	return len(this.slice)
}
