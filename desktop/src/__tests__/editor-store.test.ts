import { describe, it, expect, beforeEach } from "vitest";
import { useEditorStore, getLanguageFromPath } from "../stores/editor-store";

describe("EditorStore", () => {
  beforeEach(() => {
    useEditorStore.setState({
      openFiles: [],
      activeFilePath: null,
      projectRoot: null,
    });
  });

  it("opens a file and sets it active", () => {
    useEditorStore.getState().openFile({
      path: "/test/foo.ts",
      name: "foo.ts",
      content: "const x = 1;",
      language: "typescript",
      isDirty: false,
    });
    expect(useEditorStore.getState().openFiles).toHaveLength(1);
    expect(useEditorStore.getState().activeFilePath).toBe("/test/foo.ts");
  });

  it("does not duplicate open files", () => {
    const file = { path: "/test/foo.ts", name: "foo.ts", content: "", language: "typescript", isDirty: false };
    useEditorStore.getState().openFile(file);
    useEditorStore.getState().openFile(file);
    expect(useEditorStore.getState().openFiles).toHaveLength(1);
  });

  it("closes a file and updates active", () => {
    const store = useEditorStore.getState();
    store.openFile({ path: "/a.ts", name: "a.ts", content: "", language: "typescript", isDirty: false });
    store.openFile({ path: "/b.ts", name: "b.ts", content: "", language: "typescript", isDirty: false });
    store.closeFile("/b.ts");
    expect(useEditorStore.getState().openFiles).toHaveLength(1);
    expect(useEditorStore.getState().activeFilePath).toBe("/a.ts");
  });

  it("marks file as dirty on content update", () => {
    useEditorStore.getState().openFile({
      path: "/test.ts", name: "test.ts", content: "old", language: "typescript", isDirty: false,
    });
    useEditorStore.getState().updateFileContent("/test.ts", "new");
    expect(useEditorStore.getState().openFiles[0].isDirty).toBe(true);
  });
});

describe("getLanguageFromPath", () => {
  it("maps .ts to typescript", () => {
    expect(getLanguageFromPath("foo.ts")).toBe("typescript");
  });
  it("maps .py to python", () => {
    expect(getLanguageFromPath("script.py")).toBe("python");
  });
  it("returns plaintext for unknown extensions", () => {
    expect(getLanguageFromPath("file.xyz")).toBe("plaintext");
  });
});
