# 1.0.0 (2025-07-18)


### Bug Fixes

* abort on token error ([72e130a](https://github.com/DanArmor/vtuber-go/commit/72e130a3a8b7021feb8c1afc6d99e948f098d51f))
* auth_date zeroing on nil error fixed ([08a1a74](https://github.com/DanArmor/vtuber-go/commit/08a1a74bf8f11b0931decca25a2fd7f3f8031c6a))
* authDate fix ([b09b6eb](https://github.com/DanArmor/vtuber-go/commit/b09b6eb39cc86c72ce450f7ae92158803075ee91))
* cors fix ([a15bc8f](https://github.com/DanArmor/vtuber-go/commit/a15bc8fa19c23f1dc1499ddffdf64021e735e488))
* counter increment ([512b4fa](https://github.com/DanArmor/vtuber-go/commit/512b4fac244d17c4e982eaded6dd4261b2e7dfeb))
* docker compose proper volume mount ([39b7182](https://github.com/DanArmor/vtuber-go/commit/39b71828bf6f97d5a039d8c8a2ec550a378119a6))
* fields for video entity ([6ca6b1d](https://github.com/DanArmor/vtuber-go/commit/6ca6b1d5c2826e0dfe1de8339c1ba80d7b7ab142))
* page number validation ([47b71fa](https://github.com/DanArmor/vtuber-go/commit/47b71fa839109c341f3ac9dc1c6d16add50b0d9c))
* payload extract from context instead of initData ([c467cd5](https://github.com/DanArmor/vtuber-go/commit/c467cd5fc9c0c5df8f53d627b478affe47b7b470))
* remove migrate container from docker compose ([4fe3dba](https://github.com/DanArmor/vtuber-go/commit/4fe3dbadd439903a635316f5cc0e930927d68c9c))
* unauthorized status on token error ([6fcc66e](https://github.com/DanArmor/vtuber-go/commit/6fcc66ed1d0b5689dc1e506862bdd200b71716eb))


### Features

* added admin endpoint to add vtubers ([9b772d0](https://github.com/DanArmor/vtuber-go/commit/9b772d05a386fa1a7fd286825bab9bc23c5a9d07))
* added jwt maker to service ([a5d4683](https://github.com/DanArmor/vtuber-go/commit/a5d4683e0cc49dc6b37ce9fa3c24f382ed607b42))
* added vtuber decsription ([ef6fa39](https://github.com/DanArmor/vtuber-go/commit/ef6fa39b0c0ad14de698f24db01d65aff01250c8))
* admin token verify to the endpoint ([bbeba24](https://github.com/DanArmor/vtuber-go/commit/bbeba240592db8f059ff2e7cd93abb49d9a15c22))
* auth package ([d416537](https://github.com/DanArmor/vtuber-go/commit/d4165375ea377aee5344a2140bf6811f27b679a4))
* auth payload have userId now ([0b5b86b](https://github.com/DanArmor/vtuber-go/commit/0b5b86b1310931bb2c92f42a17ed64588d64b7d8))
* banner url field was added to vtuber ([18e4f2a](https://github.com/DanArmor/vtuber-go/commit/18e4f2a9afda3e3f3ea4347ace86bd350741d0f4))
* basic notify ([f0d91db](https://github.com/DanArmor/vtuber-go/commit/f0d91dbcde76ba3be86909043a046c66354dda12))
* breaking changes to migrations ([e3958b0](https://github.com/DanArmor/vtuber-go/commit/e3958b0d13b85420197936bdf055b8bcf727a858))
* case insensitive search ([242cdb2](https://github.com/DanArmor/vtuber-go/commit/242cdb2e7327bb92f58fa73e42aa6f69e4cd531b))
* changed from pages to offset ([6d5ad60](https://github.com/DanArmor/vtuber-go/commit/6d5ad60990a9a68c655dd3e331d7257636cdca37))
* debug flag ([1b3a747](https://github.com/DanArmor/vtuber-go/commit/1b3a7472934b3c7550f2678250c6ea7a51eaf0b9))
* deps for auth ([b1f487f](https://github.com/DanArmor/vtuber-go/commit/b1f487f378b89cf7f31e0eefb908de72a3cc59df))
* don't show twitch if no link ([8d2ad20](https://github.com/DanArmor/vtuber-go/commit/8d2ad2097bcc36268b97fcbe06c41525434842f2))
* dont mark as processed if no notify ([f72648f](https://github.com/DanArmor/vtuber-go/commit/f72648fe9b2de4014106d9d7d686893e7f2eddc8))
* embed migrations ([5cb347f](https://github.com/DanArmor/vtuber-go/commit/5cb347f347a3f74e990d423e787776eecf7e186a))
* first version of scheduler ([9eb8b42](https://github.com/DanArmor/vtuber-go/commit/9eb8b4284f40f438b14a0a3f459bdcc8d3ee2b87))
* info about utc shift ([e66d9d2](https://github.com/DanArmor/vtuber-go/commit/e66d9d2c5f502af90d594c7d61e7af18b99433cc))
* init ([c314b5f](https://github.com/DanArmor/vtuber-go/commit/c314b5fa5d2c2ab030b2db1f1aefda005b0214e5))
* init and build scripts ([aff7845](https://github.com/DanArmor/vtuber-go/commit/aff7845e2640b0976a4615779b9d76d6cdb5c93b))
* notifications system ([9ec0dd3](https://github.com/DanArmor/vtuber-go/commit/9ec0dd3c7406240c865591f968489d7e5f3cd38e))
* pagination ([1ee6d46](https://github.com/DanArmor/vtuber-go/commit/1ee6d46add32554d383473c2dc77017ea6c67c7f))
* proper concurrency protection for queue workers ([b3c5617](https://github.com/DanArmor/vtuber-go/commit/b3c56178ce2a590356aa77ffcb34f5311201455e))
* proper docker setup ([279aa50](https://github.com/DanArmor/vtuber-go/commit/279aa50136e9f97e6b233bdf383efb77ef13e9d5))
* response format ([acb5e4b](https://github.com/DanArmor/vtuber-go/commit/acb5e4b815797511849431e2e88fb8f8132107c3))
* return user id if connected to vtuber ([85a5bec](https://github.com/DanArmor/vtuber-go/commit/85a5beccb094a113fc2ca4e958e4e860b81423a0))
* rework for tasks queue ([a738fea](https://github.com/DanArmor/vtuber-go/commit/a738fea49b55f5c120a4474f90c3a7b23b82e1bf))
* reworked worker, added container for it, limited logs ([707ac3b](https://github.com/DanArmor/vtuber-go/commit/707ac3b33207e7f74fffbacf79fdf4d489d24f20))
* selecte endpoint ([8d47ab9](https://github.com/DanArmor/vtuber-go/commit/8d47ab92fd7ee027b07260c0d50790c111a6874e))
* start command and readme ([4b684b8](https://github.com/DanArmor/vtuber-go/commit/4b684b8fdb8dc13563355b3d0588d259ff783647))
* start for queue ([d165473](https://github.com/DanArmor/vtuber-go/commit/d16547360454c2f8f99bfc36f9a67960b1ebb349))
* start of admin api ([9b339b5](https://github.com/DanArmor/vtuber-go/commit/9b339b5b5ac9536c1351e3fe77dc4537192f8f2d))
* tg_init_data is separate endpoint now ([4b015fd](https://github.com/DanArmor/vtuber-go/commit/4b015fd24da7378ec5c78a183f1c34a44d91f463))
* token protection for api ([c78e5fe](https://github.com/DanArmor/vtuber-go/commit/c78e5fe9e77902118c8514a082aa0155cf5fb010))
* types/vtuber json fields changed according to holodex ([095f1ed](https://github.com/DanArmor/vtuber-go/commit/095f1ed3fc4828fb7f2700c2de09a7b9bd73e60b))
* user tg id is unique now ([8d18a18](https://github.com/DanArmor/vtuber-go/commit/8d18a181ab53b6f898d0cc01f86752d80c9e90f4))
* wave and company in search response ([5fbe4f2](https://github.com/DanArmor/vtuber-go/commit/5fbe4f26c8a4a6d614059f5f714e91be7ef93493))
* youtube channel id added ([47b938e](https://github.com/DanArmor/vtuber-go/commit/47b938e5446aa5abfafe10589fa94083c48f44bd))
* zap and proper router setup ([b2a786e](https://github.com/DanArmor/vtuber-go/commit/b2a786e67eb242500d5ff93fd079ea7a2f63d50a))


### BREAKING CHANGES

* move to scheduler that spawns work for workers in queue
* migrate from atlas migrate tool to go-migrate and embeded migrations
