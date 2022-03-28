# Copyright 2017 The Kubernetes Authors. All rights reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

<<<<<<< HEAD
FROM golang:1.17.5
=======
FROM golang:1.16
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
LABEL maintainer="Marcin Wielgus <mwielgus@google.com>"

ENV GOPATH /gopath/
ENV PATH $GOPATH/bin:$PATH
ENV GO111MODULE auto

RUN apt-get update && apt-get --yes install libseccomp-dev
RUN go version
RUN go get github.com/tools/godep
RUN godep version
CMD ["/bin/bash"]
