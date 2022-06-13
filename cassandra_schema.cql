CREATE KEYSPACE IF NOT EXISTS image_suggestions WITH replication = {'class': 'NetworkTopologyStrategy', 'eqiad': '3'};

-- The recommended images dataset
CREATE TABLE IF NOT EXISTS image_suggestions.suggestions (
    wiki          text,       -- Wikimedia project
    page_id       int,        -- MediaWiki page ID attribute
    id            timeuuid,   -- Globally unique, but also a timestamp; Unique to an algorithm run
    image         text,       -- Image being recommended
    origin_wiki   text,       -- Where the recommended image resides
    confidence    float,      -- Strength of recommendation; Value in the range 0.0-1.0
    found_on      set<text>,  -- Other wikis that use the image
    kind          set<text>,  -- ??
    page_rev      int,        -- Revision of page_id at time of recommendation (informational)
    PRIMARY KEY((wiki, page_id), id, image)
);

-- A record of user feedback, replicating whatever attributes of the corresponding
-- suggestions that is necessary.
CREATE TABLE IF NOT EXISTS image_suggestions.feedback(
    wiki          text,       -- Corresponds to suggestions.wiki
    page_id       int,        -- Corresponds to suggestions.page_id
    filename      text,       -- Corresponds to suggestions.image
    id            timeuuid,   -- ID (& timestamp) of feedback
    origin_wiki   text,       -- Corresponds to suggestions.origin_wiki
    user_id       int,        -- User who submitted feedback
    is_accepted   boolean,    -- True if feedback indicates acceptance
    is_rejected   boolean,    -- True if feedback indicates rejection
    rejection_reason  text,   -- User-submitted comment for a rejection
    PRIMARY KEY((wiki, page_id), filename, id)
);

-- Page ID/page title mapping.

-- NOTE: This table is a duplication of a relationship that MediaWiki is canonical
-- for.  It is maintained here for convenience, with the understanding that it is
-- not trustworthy (it should not be considered a source of truth).
CREATE TABLE IF NOT EXISTS image_suggestions.title_cache (
    wiki          text,       -- Wikimedia project
    page_id       int,        -- MediaWiki page ID attribute
    page_rev      int,        -- Revision of page_id
    title         text,       -- Title of page at corresponding page_rev
    PRIMARY KEY((wiki, title))
);

-- Values of the P31 property for the Wikidata item that corresponds with the page.

-- NOTE: This table is a duplication of a relationship that MediaWiki is canonical
-- for.  It is maintained here for convenience, with the understanding that it is
-- not trustworthy (it should not be considered a source of truth).
CREATE TABLE IF NOT EXISTS image_suggestions.instanceof_cache (
    wiki          text,       -- Wikimedia project
    page_id       int,        -- MediaWiki page ID attribute
    page_rev      int,        -- Revision of page_id (FIXME: shouldn't this be the Wikidata page_rev?)
    instance_of   set<text>,  -- P31 property values
    PRIMARY KEY((wiki, page_id))
);
