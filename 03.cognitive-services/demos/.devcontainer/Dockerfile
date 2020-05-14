#-------------------------------------------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See https://go.microsoft.com/fwlink/?linkid=2090316 for license information.
#-------------------------------------------------------------------------------------------------------------

FROM golang:1.14

ENV DEBIAN_FRONTEND=noninteractive
# Avoid warnings by switching to noninteractive

# This Dockerfile adds a non-root user with sudo access. Use the "remoteUser"
# property in devcontainer.json to use it. On Linux, the container user's GID/UIDs
# will be updated to match your local UID/GID (when using the dockerFile property).
# See https://aka.ms/vscode-remote/containers/non-root-user for details.
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG DOCKER_GUID=966

ARG GHCLI_VERSION=0.8.0

# Configure apt, install packages and tools
RUN apt-get update \
    && apt-get -y install --no-install-recommends apt-utils dialog 2>&1 \
    #
    # Verify git, process tools, lsb-release (common in install instructions for CLIs) installed
    && apt-get -y install git openssh-client less iproute2 procps lsb-release \
    #
    #
    # Install Go tools w/module support
    && mkdir -p /tmp/gotools \
    && cd /tmp/gotools \
    && GO111MODULE=on go get -v golang.org/x/tools/gopls@latest 2>&1 \
    && GO111MODULE=on go get -v \
        honnef.co/go/tools/...@latest \
        golang.org/x/tools/cmd/gorename@latest \
        golang.org/x/tools/cmd/goimports@latest \
        golang.org/x/tools/cmd/guru@latest \
        golang.org/x/lint/golint@latest \
        github.com/mdempsky/gocode@latest \
        github.com/cweill/gotests/...@latest \
        github.com/haya14busa/goplay/cmd/goplay@latest \
        github.com/sqs/goreturns@latest \
        github.com/josharian/impl@latest \
        github.com/davidrjenni/reftools/cmd/fillstruct@latest \
        github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest  \
        github.com/ramya-rao-a/go-outline@latest  \
        github.com/acroca/go-symbols@latest  \
        github.com/godoctor/godoctor@latest  \
        github.com/rogpeppe/godef@latest  \
        github.com/zmb3/gogetdoc@latest \
        github.com/fatih/gomodifytags@latest  \
        github.com/mgechev/revive@latest  \
        github.com/go-delve/delve/cmd/dlv@latest 2>&1 \
    #
    # Install Go tools w/o module support
    && go get -v github.com/alecthomas/gometalinter 2>&1 \
    #
    # Install gocode-gomod
    && go get -x -d github.com/stamblerre/gocode 2>&1 \
    && go build -o gocode-gomod github.com/stamblerre/gocode \
    && mv gocode-gomod $GOPATH/bin/ \
    #
    # Install golangci-lint
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin 2>&1 \
    #
    # Create a non-root user to use if preferred - see https://aka.ms/vscode-remote/containers/non-root-user.
    && groupadd --gid $USER_GID $USERNAME \
    && useradd -s /bin/bash --uid $USER_UID --gid $USER_GID -m $USERNAME \
    # [Optional] Add sudo support
    && apt-get install -y sudo \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME \
    #
    # zsh, oh-my-zsh and tools
    && apt-get install -y vim less zsh tree \
    && mkdir -p /tmp/zsh \
    && cd /tmp/zsh \
    && su ${USERNAME} sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)" \
    && echo 'alias ff="find . -iname " >> /home/${USERNAME}/.zshrc' \
    && echo 'alias ls="ls -h --color=auto " >> /home/${USERNAME}/.zshrc' \
    && echo 'alias ll="ls -l " >> /home/${USERNAME}/.zshrc' \
    && echo 'alias l="ls -l " >> /home/${USERNAME}/.zshrc' \
    && echo 'alias la="ls -al " >> /home/${USERNAME}/.zshrc' \
    && echo 'alias l.="ls -al " >> /home/${USERNAME}/.zshrc' \
    && sed -i 's/^plugins=\(.*\)$/plugins=(git kubectl minikube helm github golang docker)/' /home/${USERNAME}/.zshrc \
    #
    # Install Starship
    && mkdir -p /tmp/starship \
    && cd /tmp/starship \
    && curl -ssSL https://starship.rs/install.sh -o install.sh \
    && chmod +x install.sh \
    && ./install.sh -y \
    && echo "# ~/.zshrc" >> /home/${USERNAME}/.zshrc \
    && echo "eval \"$(starship init zsh)\"" >> /home/${USERNAME}/.zshrc \
    #
    # GitHub CLI
    && mkdir -p /tmp/gh \
    && cd /tmp/gh \
    && curl -sLO  https://github.com/cli/cli/releases/download/v${GHCLI_VERSION}/gh_${GHCLI_VERSION}_linux_amd64.deb \
    && apt install ./gh_${GHCLI_VERSION}_linux_amd64.deb \
    && echo "eval $(gh completion --shell zsh)" >> /home/${USERNAME}/.zshrc \
    #
    # Permit VS Code find the installed extensions after container rebuild
    && mkdir -p /home/$USERNAME/.vscode-server/extensions \
        /home/$USERNAME/.vscode-server-insiders/extensions \
    && chown -R $USERNAME \
        /home/$USERNAME/.vscode-server \
        /home/$USERNAME/.vscode-server-insiders  \
    #
    # Install kubectl
    && curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl \
    && chmod +x kubectl \
    && mv ./kubectl /usr/bin/kubectl \
    #
    ## kubens/kubectx: https://github.com/ahmetb/kubectx
    && apt install -y kubectx \
    #
    ## kubectl alias
    && curl -sSL https://raw.githubusercontent.com/ahmetb/kubectl-alias/master/.kubectl_aliases -o /home/${USERNAME}/.kubectl_aliases \
    && echo "# Kubectl aliases from https://github.com/ahmetb/kubectl-aliases" >> /home/${USERNAME}/.zshrc \
    && echo "[ -f ~/.kubectl_aliases ] && source ~/.kubectl_aliases" >> /home/${USERNAME}/.zshrc \
    && echo "# watch workaround for kubectl shortcuts" >> /home/${USERNAME}/.zshrc \
    && echo "alias watch=\"watch \"" >> /home/${USERNAME}/.zshrc \
    && chown ${USER_UID}:${USER_GID} /home/${USERNAME}/.kubectl_aliases \
    #
    # Add write permission for /go/pkg
    && chown -R ${USERNAME}:${USERNAME} /go \
    #
    # Install Docker CE CLI
    && apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common lsb-release \
    && curl -fsSL https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg | (OUT=$(apt-key add - 2>&1) || echo $OUT) \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]') $(lsb_release -cs) stable" \
    && apt-get update \
    && apt-get install -y docker-ce-cli \
    && groupadd -g ${DOCKER_GUID} docker \
    && usermod -aG docker ${USERNAME} \
    #
    # Install Docker Compose
    && curl -sSL "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
    && chmod a+x /usr/local/bin/docker-compose \
    && echo ${DOCKER_GUID} > /home/${USERNAME}/dockerguid \
    && chown ${USER_UID}:${USER_GID} /home/${USERNAME}/dockerguid \
    #
    # Helm
    && curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash \
    # 
    # Install Azure Client
    && curl -sL https://aka.ms/InstallAzureCLIDeb | bash \
    && echo "# Azure cli completion " >> /home/${USERNAME}/.zshrc \
    && echo "autoload -U +X bashcompinit && bashcompinit" >> /home/${USERNAME}/.zshrc \
    && echo "source /etc/bash_completion.d/azure-cli" >> /home/${USERNAME}/.zshrc \
    #
    ## Initialize user .config dir
    && mkdir -p /home/${USERNAME}/.config \
    && chown ${USER_UID}:${USER_GID} /home/${USERNAME}/.config \
    #
    # Clean up
    && cd / \
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/* /go/src /tmp/* 

# Update this to "on" or "off" as appropriate
ENV GO111MODULE=auto

# Switch back to dialog for any ad-hoc use of apt-get
ENV DEBIAN_FRONTEND=dialog
