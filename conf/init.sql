--  This file is part of the eliona project.
--  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
--  ______ _ _
-- |  ____| (_)
-- | |__  | |_  ___  _ __   __ _
-- |  __| | | |/ _ \| '_ \ / _` |
-- | |____| | | (_) | | | | (_| |
-- |______|_|_|\___/|_| |_|\__,_|
--
--  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
--  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
--  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
--  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
--  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

create schema if not exists microsoft_365;

-- Should be editable by eliona frontend.
create table if not exists microsoft_365.configuration
(
	id               bigserial primary key,
	client_id        text not null,
	client_secret    text not null,
	tenant_id        text not null,
	username         text not null,
	password         text not null,
	for_eliona       boolean not null default false,
	for_proxy        boolean not null default false,
	refresh_interval integer not null default 60,
	request_timeout  integer not null default 120,
	asset_filter     json,
	active           boolean default false,
	enable           boolean default false,
	project_ids      text[]
);

create table if not exists microsoft_365.asset
(
	id               bigserial primary key,
	configuration_id bigserial not null references microsoft_365.configuration(id) ON DELETE CASCADE,
	project_id       text      not null,
	global_asset_id  text      not null,
	asset_id         integer,
	email            text      not null
);

-- Makes the new objects available for all other init steps
commit;
