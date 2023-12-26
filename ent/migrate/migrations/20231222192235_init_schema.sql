-- Create "users" table
CREATE TABLE "users" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "tg_id" bigint NOT NULL, "first_name" character varying NOT NULL, "last_name" character varying NULL, "username" character varying NULL, "language_code" character varying NULL, "photo_url" character varying NULL, PRIMARY KEY ("id"));
-- Create "orgs" table
CREATE TABLE "orgs" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "name" character varying NOT NULL, PRIMARY KEY ("id"));
-- Create "waves" table
CREATE TABLE "waves" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "name" character varying NOT NULL, "org_waves" bigint NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "waves_orgs_waves" FOREIGN KEY ("org_waves") REFERENCES "orgs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create "vtubers" table
CREATE TABLE "vtubers" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "channel_name" character varying NOT NULL, "english_name" character varying NULL, "photo_url" character varying NULL, "twitter" character varying NULL, "video_count" bigint NULL, "subscriber_count" bigint NULL, "clip_count" bigint NULL, "top_topics" jsonb NULL, "inactive" boolean NOT NULL, "twitch" character varying NULL, "wave_vtubers" bigint NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "vtubers_waves_vtubers" FOREIGN KEY ("wave_vtubers") REFERENCES "waves" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create "user_vtubers" table
CREATE TABLE "user_vtubers" ("user_id" bigint NOT NULL, "vtuber_id" bigint NOT NULL, PRIMARY KEY ("user_id", "vtuber_id"), CONSTRAINT "user_vtubers_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "user_vtubers_vtuber_id" FOREIGN KEY ("vtuber_id") REFERENCES "vtubers" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);