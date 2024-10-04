import { useMeta } from "../state/use-meta";
import HomeMeta from "./chunks/HomeMeta";

export default function Home() {
  const { metadata } = useMeta();

  return (
    <>
      <HomeMeta {...metadata} />
    </>
  )
}
