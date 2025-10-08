# kubelocal - Spécification Technique et Fonctionnelle

## 1. Vue d'ensemble

### 1.1 Objectif
`kubelocal` est un outil d'installation interactif pour environnements Kubernetes locaux, conçu pour l'apprentissage et l'expérimentation. Il cible les environnements temporaires (VMs, WSL) et non la production.

### 1.2 Scope
- **Utilisateurs cibles** : Débutants Kubernetes, développeurs en phase d'apprentissage
- **Environnements** : VMs, WSL2, machines de développement Linux
- **Usage** : Installation one-shot, environnements éphémères (pas d'uninstall nécessaire)

### 1.3 Plateformes
- **Actuel** : Linux AMD64 (Debian-based)
- **Futur** : Linux ARM64

---

## 2. Spécification Fonctionnelle

### 2.1 Flow Utilisateur

```
1. Lancement → kubelocal
2. Affichage du welcome banner
3. Sélection interactive :
   - Solution Kubernetes (Kind / Minikube / K3s)
   - Outils de manifest management (Helm / Kustomize) [multi-select]
   - [FUTUR] Service mesh (Istio)
4. Résumé de la configuration
5. Installation automatique avec feedback visuel (spinner)
6. Affichage des instructions de démarrage personnalisées
```

### 2.2 Composants Installés

#### 2.2.1 Toujours Installés
- **kubectl** : Latest version via dl.k8s.io/release/stable.txt
- **jq** : Via apt (warning si échec, non-critique)
- **Aliases bash** : `k=kubectl`, `logs` pour jq (warning si échec)
  - Autocomplétion kubectl activée via `source <(kubectl completion bash)`
  - Autocomplétion étendue à l'alias `k` via `complete -o default -F __start_kubectl k`
  - Futur : `cheat` pour afficher cheat sheets avec glow

#### 2.2.2 Conditionnels
- **Docker** : Si Kind ou Minikube sélectionné
  - Configuration permissions automatique (usermod -aG docker)
  - Workarounds temporaires : sg docker, chmod 666 docker.sock
- **Kind** : v0.30.0 (hardcodé, à améliorer)
- **Minikube** : Latest release
- **K3s** : Latest via get.k3s.io (téléchargé puis exécuté en deux étapes)
  - Inclut configuration kubectl automatique (~/.kube/config)
- **Helm** : Si sélectionné, via script officiel
- **Kustomize** : Si sélectionné, via script officiel

#### 2.2.3 Futur
- **Istio** : istioctl + installation sur cluster
- **Support ARM64**
- **Podman** (mentionné dans README mais pas implémenté)
- **Golang** : Installation du langage Go
- **Glow** : Pour affichage markdown des cheat sheets

### 2.3 Gestion des Erreurs

#### Erreurs Critiques (exit 1)
- Échec installation composant principal (Docker, Kind, Minikube, K3s, kubectl)
- Échec form interactif
- Échec logger

#### Warnings (continue)
- Échec jq installation
- Échec setup aliases
- Échec Istio (si implémenté)

**Problème identifié** : Distinction entre erreurs critiques/warnings à documenter et potentiellement revoir.

---

## 3. Spécification Technique

### 3.1 Architecture Actuelle

```
kubelocal/
├── main.go                    # Entry point, orchestration, UI
├── tools/                     # Logique installation (à renommer packages/ ?)
│   ├── common.go              # kubectl, jq, aliases
│   ├── docker.go              # Docker + permissions
│   ├── kind.go
│   ├── minikube.go
│   ├── k3s.go
│   ├── helm.go
│   ├── kustomize.go
│   └── istio.go               # Code présent mais désactivé
└── internal/
    ├── logger/                # Logging structuré (zap)
    └── command/               # Wrapper exec.Command

```

#### Points d'amélioration
1. **Naming** : `tools/` → `packages/` plus idiomatique Go
2. **Séparation** : `internal/` pour code réutilisable interne, `packages/` pour logique métier
3. **Structure** : Architecture plate OK pour MVP, à évaluer si croissance

### 3.2 Logging

#### État Actuel - PROBLÉMATIQUE
- **Double output** : `fmt.Printf` (user) + `log.Info` (fichier)
- **Propagation erreurs** : Erreurs wrappées plusieurs niveaux (`fmt.Errorf` → `fmt.Errorf` → log.Error)
- **Incohérence** : Certains messages en console seulement, d'autres loggés

#### Exemple du problème
```go
// tools/docker.go:36
if err := setupDockerPermissions(log); err != nil {
    return fmt.Errorf("failed to setup Docker permissions: %w", err)
}
// → main.go:86 handleInstallationError log cette erreur
// → Double/triple logging de la même erreur
```

#### Améliorations Nécessaires
1. **Séparation claire** :
   - Console : Feedback utilisateur simple, progressif
   - Logs : Détails techniques, debugging
2. **Niveaux de verbosité** : --verbose, --debug flags
3. **Format** : JSON logs optionnel pour parsing
4. **Rotation** : Non nécessaire (environnements éphémères)

### 3.3 Flux d'Exécution

```
main()
  ├─> showWelcome()
  ├─> runSetup(log)
      ├─> Form interactif (huh library)
      ├─> Détermination dépendances (Docker needed?)
      ├─> installComponent() pour chaque outil
      │   ├─> runSpinner() (goroutine)
      │   └─> tools.Install*()
      │       └─> command.Run()
      └─> showGettingStarted()
```

### 3.4 Composants Techniques

#### 3.4.1 Logger (`internal/logger`)
- **Librairie** : uber/zap
- **Outputs** : File (~/.local/share/kubelocal/logs/install.log)
- **Niveaux** : Info, Debug, Warn, Error
- **Format** : JSON structuré

#### 3.4.2 Command Runner (`internal/command`)
- Wrapper exec.Command
- Logging automatique stdin/stdout/stderr
- Propagation erreurs

#### 3.4.3 UI
- **Librairie** : charmbracelet/huh (forms TUI)
- **Spinner** : Custom implementation (chars: | / - \)
- **Problème identifié** : `\r` efface ligne, peut perdre messages si échec

### 3.5 Gestion Versions

#### Actuel
- **kubectl** : Latest via API Kubernetes
- **Kind** : v0.30.0 hardcodé
- **Autres** : Latest via scripts officiels

#### Cible
- **Aucune version hardcodée**
- Tout en latest/stable dynamique
- Possible paramétrage futur (`--kubectl-version`)

---

## 4. Tests (À Implémenter)

### 4.1 Tests Unitaires (Priorité 1)

#### Packages à tester
```go
// tools/ → packages/
TestInstallKubectl()
TestInstallDocker()
TestSetupDockerPermissions()
TestFindIstioDirectory()
// etc.

// internal/command
TestRunCommand()
TestRunWithOutput()

// main.go helpers
TestDetermineDockerNeed()
```

#### Stratégie
- **Mock commands** : Interface pour exec.Command
- **Filesystem mock** : afero ou similar
- **Isolation** : Pas d'installation réelle
- **Coverage** : Minimum 70%

### 4.2 Tests d'Intégration (Futur)
- Container Docker pour environnement test
- Installation réelle dans environnement isolé

### 4.3 Tests E2E (Futur)
- VM éphémère (Vagrant/Terraform)
- Flow complet : install → vérification → cleanup

### 4.4 Architecture pour Tests

**Problème actuel** : Functions directement liées à exec.Command, difficile à tester

**Solution** :
```go
// interface commandRunner
type CommandRunner interface {
    Run(cmd string, args ...string) error
    RunWithOutput(cmd string, args ...string) (string, error)
}

// Production: realCommandRunner
// Tests: mockCommandRunner
```

---

## 5. Configuration & CLI

### 5.1 Flags Actuels
- `--version` : Affiche version

### 5.2 Flags Futurs
- `--verbose` : Logs détaillés en console
- `--debug` : Debug mode (all logs)
- `--non-interactive` : Config via flags (CI/CD)
- `--log-format json` : Format logs

---

## 6. Shell Support

### Actuel
- **bash** : Configuration complète (voir section 2.2.1)
  - Message affiché : utilisateur doit faire `source ~/.bashrc` ou redémarrer shell
- **Limitation** : Source automatique impossible (subprocess ne peut modifier shell parent)

### Futur
- **zsh** : Support prévu
- **fish** : À évaluer
- **Détection automatique** : $SHELL

---

## 7. Dépendances

### Runtime
- `curl` : Téléchargements
- `sudo` : Installations système
- `apt` : Package manager (Debian-based)

### Futures (à installer)
- `glow` : Markdown renderer pour cheat sheets
- `golang` : Langage Go (optionnel ou toujours installé, à décider)

### Go Modules
```
github.com/charmbracelet/huh  # TUI forms
go.uber.org/zap               # Logging
```

---

## 8. Sécurité & Permissions

### 8.1 Docker Permissions - Workaround Temporaire

**Problème** : Après `usermod -aG docker`, groupe pas actif dans session actuelle

**Solutions tentées** :
1. `sg docker` : Active groupe temporairement
2. `chmod 666 /var/run/docker.sock` : Permission socket (⚠️ insecure)

**Status** : À tester, solution 2 non recommandée pour prod (OK pour learning env)

**Recommandation** :
- Informer utilisateur de logout/login
- Tester `newgrp docker`
- Documenter limitations

### 8.2 Sudo
- Requis pour installations système
- Pas de sudo cache management
- User doit avoir droits sudo

---

## 9. Outputs & Logs

### 9.1 Console (stdout)
```
✅ Success messages
❌ Error messages
⚠️  Warning messages
🚀 Getting started instructions
```

### 9.2 Logs (file)
- **Path** : `~/.local/share/kubelocal/logs/install.log`
- **Format** : JSON structuré (zap)
- **Alias** : `logs` command pour jq parsing

### 9.3 Problème à Résoudre
Clarifier responsabilités :
- Quoi logger vs quoi afficher
- Éviter double/triple logging erreurs
- Cohérence messages

---

## 10. Kubernetes-specific

### 10.1 Kind
- Docker requis
- Cluster creation: Manuel post-install
- Multi-cluster support

### 10.2 Minikube  
- Docker requis
- Drivers: docker (default)
- Dashboard inclus

### 10.3 K3s
- No Docker needed (containerd intégré)
- Systemd service
- Auto-configure kubectl (~/.kube/config)
- Permissions: sudo needed

---

## 11. Cheat Sheets (Feature Future)

### 11.1 Concept
Centraliser les cheat sheets de tous les outils installés dans un fichier markdown unique, facilement accessible via un alias.

### 11.2 Implémentation Prévue

#### Fichier cheat sheet
- **Localisation** : `~/.local/share/kubelocal/cheatsheet.md`
- **Contenu** : Commandes essentielles pour chaque outil installé
- **Structure** :
  ```markdown
  # Kubelocal Cheat Sheet
  
  ## kubectl / k
  - Get pods: `k get pods`
  - Describe: `k describe pod <name>`
  
  ## Helm
  - Install chart: `helm install <name> <chart>`
  - List releases: `helm list`
  
  ## Kind
  - Create cluster: `kind create cluster --name <name>`
  ...
  ```

#### Génération dynamique
- Créé pendant l'installation
- Contenu adapté aux outils sélectionnés
- Sections générées uniquement pour outils installés

#### Affichage
- **Outil** : `glow` (markdown renderer TUI)
  - Installation automatique pendant setup
  - Fallback sur `cat` si glow non disponible
- **Alias** : `cheat` ou `kl-cheat`
  - Ajouté à ~/.bashrc et ~/.zshrc
  - Commande : `glow ~/.local/share/kubelocal/cheatsheet.md`

### 11.3 Contenu par Outil

#### kubectl
- Commandes CRUD de base
- Contexts et namespaces
- Logs et debug
- Common flags

#### Helm
- Repo management
- Install/upgrade/rollback
- Values overrides

#### Kustomize
- Build et apply
- Structure directories
- Overlays

#### Kind
- Cluster lifecycle
- Load images
- Multi-node clusters

#### Minikube
- Start/stop avec options
- Addons
- Tunnel et port-forward

#### K3s
- Systemctl commands
- Kubeconfig setup
- Logs journalctl

#### Istio (futur)
- Installation profiles
- Traffic management
- Observability tools

### 11.4 Maintenance
- **Updates** : Fichier statique post-install (pas de update automatique)
- **Versioning** : Cheat sheet versionné avec kubelocal
- **Personnalisation** : User peut éditer librement

---

## 12. Features Futures (Roadmap)

### Phase 1 - Tests & Refactoring
- [ ] Tests unitaires (mocking commands)
- [ ] Refactor logging (séparation console/file)
- [ ] Améliorer architecture testabilité
- [ ] Renommer tools/ → packages/

### Phase 2 - Versions & Support
- [ ] Supprimer versions hardcodées
- [ ] Support zsh
- [ ] Detection shell automatique
- [ ] ARM64 support

### Phase 3 - Features
- [ ] Istio installation
- [ ] Podman support (alternatif Docker)
- [ ] Mode non-interactif (flags)
- [ ] Verbosity levels
- [ ] Installation Golang
- [ ] Cheat sheets centralisées
  - [ ] Installation glow
  - [ ] Génération cheatsheet.md dynamique
  - [ ] Alias `cheat` avec glow

### Phase 4 - Polish
- [ ] Améliorer spinner (preserve messages)
- [ ] Better error messages
- [ ] Validation pré-installation (disk space, etc.)
- [ ] Progress bars (alternative spinner)

---

## 13. Limitations & Contraintes

### Techniques
- Linux only (pas Windows natif, WSL OK)
- AMD64 only (ARM futur)
- Debian-based (apt dependency)
- Bash dependency (zsh futur)

### Design Choices
- Pas d'uninstall (environnements éphémères)
- Latest versions (pas de version pinning user)
- Interactive uniquement (non-interactive futur)
- Single binary (pas de config files)

### Known Issues
- Docker permissions workaround temporaire
- Spinner efface messages sur erreur
- Logs incohérents (fmt vs zap)
- Versions hardcodées (Kind)

---

## 14. Build & Release

### Version
- Format: `v0.0.5` (semver)
- Embedded in binary: `-ldflags "-X main.version=..."`

### Distribution
- Website installer: `curl | sh`
- GitHub releases: tar.gz binaries
- Pas de package managers (snap/apt) pour l'instant

---

## 15. Documentation

### Existante
- README.md : User-facing
- spec.md : Ce document
- cmd.md : (à vérifier contenu)

### Manquante
- Architecture diagrams
- Contributing guide
- Testing guide
- Development setup

### Future
- Cheat sheet intégré (`~/.local/share/kubelocal/cheatsheet.md`)
  - Généré dynamiquement selon outils installés
  - Accessible via alias `cheat` + glow

---

## 16. Métriques de Succès

### Installation
- ✅ Binaire exécutable sans config
- ✅ Installation complète < 10min
- ✅ Feedback temps réel (spinner)
- ✅ Instructions post-install claires

### Qualité
- ⏳ Tests unitaires (à faire)
- ⏳ Logs cohérents (à améliorer)
- ⚠️  Gestion erreurs (partiel)
- ✅ Zero dependencies à installer manuellement

### User Experience
- ✅ Interactive TUI
- ✅ Choix clairs
- ⏳ Error messages helpful (à améliorer)
- ✅ Getting started guide

---

## 17. Questions Ouvertes

1. **Architecture internal/** : Garder ou tout mettre dans packages/ ?
2. **Command mocking** : Interface ou struct avec dependency injection ?
3. **Logging levels** : Implementer maintenant ou après tests ?
4. **Istio priority** : Quand l'activer ? Après tests ou avant ?
5. **ARM support** : Cross-compile ou native builds ?
6. **Shell detection** : Modifier ~/.bashrc ET ~/.zshrc ? Ou juste celui actif ?
7. **Cheat sheets** : Générer contenu statique ou templates Go ? Quelle profondeur de contenu ?
8. **Golang installation** : Toujours installer ou option utilisateur ? Quelle version (latest/LTS) ?
9. **Glow fallback** : Si glow install échoue, `cat` suffit ou autre renderer (bat, rich) ?

---

*Document vivant - À mettre à jour avec évolutions du projet*

