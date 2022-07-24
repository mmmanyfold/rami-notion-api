import { Client } from "@notionhq/client";
import { GetPagePropertyResponse, PageObjectResponse, PartialPageObjectResponse, QueryDatabaseResponse } from "@notionhq/client/build/src/api-endpoints";
import dotenv from "dotenv";

dotenv.config();
const databaseId = process.env.NOTION_DATABASE_ID || "";

const notion = new Client({
  auth: process.env.NOTION_TOKEN
});

/**
 * Gets projects from the database.
 *
 * @returns {Promise<Array<{ pageId: string, status: string, title: string }>>}
 */
 async function getProjectsFromNotionDatabase() {
  const pages: Array<PageObjectResponse | PartialPageObjectResponse> = []
  let cursor = undefined

  while (true) {
    const { results, next_cursor }: QueryDatabaseResponse = await notion.databases.query({
      database_id: databaseId,
      start_cursor: cursor,
    })
    pages.push(...results)
    if (!next_cursor) {
      break
    }
    cursor = next_cursor
  }
  console.log(`${pages.length} pages successfully fetched.`)

  // TODO: FIX SCHEMA
  const projects = []
  for (const page of pages) {
    const pageId = page.id

    const projectIdProperty = page.properties["ID"].id
    const projectIdPropertyItem = await getPropertyValue({
      pageId,
      propertyId: projectIdProperty,
    })
    const projectId = projectIdPropertyItem.select
      ? projectIdPropertyItem.select.name
      : "No Project ID"

    const titlePropertyId = page.properties["Title"].id
    const titlePropertyItems = await getPropertyValue({
      pageId,
      propertyId: titlePropertyId,
    })
    const title = titlePropertyItems
      .map(propertyItem => propertyItem.title.plain_text)
      .join("")

    projects.push({ pageId, projectId, title })
  }
  return projects
}

async function main() {

  getProjectsFromNotionDatabase();

  const response = await notion.databases.query({
    database_id: databaseId,
  });

  // const sortedResults = response.results.sort((a, b) => (a.properties.ID > b.properties.ID) ? 1 : -1)

  console.log("Got response:", response);
}

main()
  .then(() => process.exit(0))
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });

/**
 * If property is paginated, returns an array of property items.
 *
 * Otherwise, it will return a single property item.
 */
 async function getPropertyValue({ pageId, propertyId }: { pageId: string, propertyId: string }) {
  const propertyItem: GetPagePropertyResponse = await notion.pages.properties.retrieve({
    page_id: pageId,
    property_id: propertyId,
  })
  if (propertyItem.object === "property_item") {
    return propertyItem
  }

  // Property is paginated.
  let nextCursor = propertyItem.next_cursor
  const results = propertyItem.results

  while (nextCursor !== null) {
    const propertyItem: GetPagePropertyResponse = await notion.pages.properties.retrieve({
      page_id: pageId,
      property_id: propertyId,
      start_cursor: nextCursor,
    })

    // TODO: FIX SCHEMA
    nextCursor = propertyItem.next_cursor
    results.push(...propertyItem.results)
  }

  return results
}
