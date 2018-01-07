# go-linode
Go client for Linode REST v4 API

**Not yet ready for production usage**

# API Support

## Linodes

- `/linode/instances`
  - [ ] `GET`
  - [ ] `POST`
- `/linode/instances/$id`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`
- `/linode/instances/$id/boot`
  - [ ] `POST`
- `/linode/instances/$id/clone`
  - [ ] `POST`
- `/linode/instances/$id/kvmify`
  - [ ] `POST`
- `/linode/instances/$id/mutate`
  - [ ] `POST`
- `/linode/instances/$id/mutate`
  - [ ] `POST`
- `/linode/instances/$id/reboot`
  - [ ] `POST`
- `/linode/instances/$id/rebuild`
  - [ ] `POST`
- `/linode/instances/$id/rescue`
  - [ ] `POST`
- `/linode/instances/$id/resize`
  - [ ] `POST`
- `/linode/instances/$id/shutdown`
  - [ ] `POST`
- `/linode/instances/$id/volumes`
  - [ ] `GET`

## Backups
- `/linode/instances/$id/backups`
  - [ ] `GET`
  - [ ] `POST`
- `/linode/instances/$id/backups/$id/restore`
  - [ ] `POST`
- `/linode/instances/$id/backups/cancel`
  - [ ] `POST`
- `/linode/instances/$id/backups/enable`
  - [ ] `POST`
  
## Configs
- `/linode/instances/$id/configs`
  - [ ] `GET`
  - [ ] `POST`
- `/linode/instances/$id/configs/$id`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

## Disks
- `/linode/instances/$id/disks`
  - [ ] `GET`
  - [ ] `POST`
- `/linode/instances/$id/disks/$id`
  - [ ] `GET` 
  - [ ] `PUT` 
  - [ ] `POST` 
  - [ ] `DELETE` 
- `/linode/instances/$id/disks/$id/imagize`
  - [ ] `POST`
- `/linode/instances/$id/disks/$id/password`
  - [ ] `POST`
- `/linode/instances/$id/disks/$id/resize`
  - [ ] `POST` 

## IPs
- `/linode/instances/$id/ips`
  - [ ] `GET`
  - [ ] `POST`
- `/linode/instances/$id/ips/$ip_address`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`
- `/linode/instances/$id/ips/sharing`
  - [ ] `POST`

## Kernels
- `/linode/kernels`
  - [ ] `GET`
- `/linode/kernels/$id`
  - [ ] `GET`

## StackScripts
- `/linode/stackscripts`
  - [ ] `GET`
  - [ ] `POST`
- `/linode/stackscripts/$id`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

## Stats
- `/linode/instances/$id/stats`
  - [ ] `GET`
- `/linode/instances/$id/stats/$year/$month`
  - [ ] `GET`

## Types
- `/linode/types`
  - [ ] `GET`
- `/linode/types/$id`
  - [ ] `GET`
