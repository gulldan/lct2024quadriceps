import Page from "./main";
import { getTask } from "@/utils/web";



const sleep = (ms = 0) => new Promise((resolve) => setTimeout(resolve, ms));

export default async function Home({ params }) {

  const data = await getTask(params.id);
  return <Page init={data} />;



}
