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
